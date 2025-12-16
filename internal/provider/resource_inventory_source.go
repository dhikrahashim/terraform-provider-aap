package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/dhikrahashim/terraform-provider-aap/internal/client"
)

var _ resource.Resource = &InventorySourceResource{}
var _ resource.ResourceWithImportState = &InventorySourceResource{}

func NewInventorySourceResource() resource.Resource {
	return &InventorySourceResource{}
}

type InventorySourceResource struct {
	client *client.Client
}

type InventorySourceResourceModel struct {
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	InventoryID        types.String `tfsdk:"inventory_id"`
	Source             types.String `tfsdk:"source"`
	SourcePath         types.String `tfsdk:"source_path"`
	SourceVars         types.String `tfsdk:"source_vars"`
	CredentialID       types.String `tfsdk:"credential_id"`
	SourceProjectID    types.String `tfsdk:"source_project_id"`
	UpdateOnLaunch     types.Bool   `tfsdk:"update_on_launch"`
	UpdateCacheTimeout types.Int64  `tfsdk:"update_cache_timeout"`
	Overwrite          types.Bool   `tfsdk:"overwrite"`
	OverwriteVars      types.Bool   `tfsdk:"overwrite_vars"`
}

func (r *InventorySourceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_inventory_source"
}

func (r *InventorySourceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Inventory source for dynamic inventory from external sources.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"inventory_id": schema.StringAttribute{
				Required: true,
			},
			"source": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Source type: scm, ec2, gce, azure_rm, vmware, satellite6, openstack, rhv, controller, file.",
			},
			"source_path": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Path to inventory file or script within project.",
			},
			"source_vars": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Source variables in YAML/JSON format.",
			},
			"credential_id": schema.StringAttribute{
				Optional: true,
			},
			"source_project_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Project containing inventory file (for scm source).",
			},
			"update_on_launch": schema.BoolAttribute{
				Optional: true,
			},
			"update_cache_timeout": schema.Int64Attribute{
				Optional: true,
			},
			"overwrite": schema.BoolAttribute{
				Optional: true,
			},
			"overwrite_vars": schema.BoolAttribute{
				Optional: true,
			},
		},
	}
}

func (r *InventorySourceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", fmt.Sprintf("Expected *client.Client, got: %T.", req.ProviderData))
		return
	}
	r.client = c
}

func (r *InventorySourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data InventorySourceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	invID, _ := strconv.Atoi(data.InventoryID.ValueString())
	credID, _ := strconv.Atoi(data.CredentialID.ValueString())
	projID, _ := strconv.Atoi(data.SourceProjectID.ValueString())

	is := &client.InventorySource{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
		Inventory:   invID,
		Source:      data.Source.ValueString(),
	}
	if !data.SourcePath.IsNull() {
		is.SourcePath = data.SourcePath.ValueString()
	}
	if !data.SourceVars.IsNull() {
		is.SourceVars = data.SourceVars.ValueString()
	}
	if credID > 0 {
		is.Credential = credID
	}
	if projID > 0 {
		is.SourceProject = projID
	}
	if !data.UpdateOnLaunch.IsNull() {
		is.UpdateOnLaunch = data.UpdateOnLaunch.ValueBool()
	}
	if !data.UpdateCacheTimeout.IsNull() {
		is.UpdateCacheTimeout = int(data.UpdateCacheTimeout.ValueInt64())
	}
	if !data.Overwrite.IsNull() {
		is.Overwrite = data.Overwrite.ValueBool()
	}
	if !data.OverwriteVars.IsNull() {
		is.OverwriteVars = data.OverwriteVars.ValueBool()
	}

	created, err := r.client.CreateInventorySource(is)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create inventory source: %s", err))
		return
	}

	data.ID = types.StringValue(strconv.Itoa(created.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InventorySourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data InventorySourceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, _ := strconv.Atoi(data.ID.ValueString())
	is, err := r.client.GetInventorySource(id)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read inventory source: %s", err))
		return
	}

	data.Name = types.StringValue(is.Name)
	data.Description = types.StringValue(is.Description)
	data.InventoryID = types.StringValue(strconv.Itoa(is.Inventory))
	data.Source = types.StringValue(is.Source)
	data.SourcePath = types.StringValue(is.SourcePath)
	data.SourceVars = types.StringValue(is.SourceVars)
	if is.Credential > 0 {
		data.CredentialID = types.StringValue(strconv.Itoa(is.Credential))
	}
	if is.SourceProject > 0 {
		data.SourceProjectID = types.StringValue(strconv.Itoa(is.SourceProject))
	}
	data.UpdateOnLaunch = types.BoolValue(is.UpdateOnLaunch)
	data.UpdateCacheTimeout = types.Int64Value(int64(is.UpdateCacheTimeout))
	data.Overwrite = types.BoolValue(is.Overwrite)
	data.OverwriteVars = types.BoolValue(is.OverwriteVars)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InventorySourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data InventorySourceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, _ := strconv.Atoi(data.ID.ValueString())
	invID, _ := strconv.Atoi(data.InventoryID.ValueString())
	credID, _ := strconv.Atoi(data.CredentialID.ValueString())
	projID, _ := strconv.Atoi(data.SourceProjectID.ValueString())

	is := &client.InventorySource{
		ID:          id,
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
		Inventory:   invID,
		Source:      data.Source.ValueString(),
	}
	if !data.SourcePath.IsNull() {
		is.SourcePath = data.SourcePath.ValueString()
	}
	if !data.SourceVars.IsNull() {
		is.SourceVars = data.SourceVars.ValueString()
	}
	if credID > 0 {
		is.Credential = credID
	}
	if projID > 0 {
		is.SourceProject = projID
	}
	if !data.UpdateOnLaunch.IsNull() {
		is.UpdateOnLaunch = data.UpdateOnLaunch.ValueBool()
	}
	if !data.UpdateCacheTimeout.IsNull() {
		is.UpdateCacheTimeout = int(data.UpdateCacheTimeout.ValueInt64())
	}
	if !data.Overwrite.IsNull() {
		is.Overwrite = data.Overwrite.ValueBool()
	}
	if !data.OverwriteVars.IsNull() {
		is.OverwriteVars = data.OverwriteVars.ValueBool()
	}

	_, err := r.client.UpdateInventorySource(is)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update inventory source: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InventorySourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data InventorySourceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, _ := strconv.Atoi(data.ID.ValueString())
	if err := r.client.DeleteInventorySource(id); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete inventory source: %s", err))
	}
}

func (r *InventorySourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
