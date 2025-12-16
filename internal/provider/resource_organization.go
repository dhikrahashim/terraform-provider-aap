package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/dhikrahashim/terraform-provider-aap/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &OrganizationResource{}
var _ resource.ResourceWithImportState = &OrganizationResource{}

func NewOrganizationResource() resource.Resource {
	return &OrganizationResource{}
}

type OrganizationResource struct {
	client *client.Client
}

type OrganizationResourceModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	MaxHosts         types.Int64  `tfsdk:"max_hosts"`
	CustomVirtualEnv types.String `tfsdk:"custom_virtualenv"`
}

func (r *OrganizationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

func (r *OrganizationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Organization resource for AAP.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Organization ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the organization.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Description of the organization.",
			},
			"max_hosts": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Maximum number of hosts allowed to be managed by this organization.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"custom_virtualenv": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Local absolute file path containing a custom Python virtualenv to use.",
			},
		},
	}
}

func (r *OrganizationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *OrganizationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data OrganizationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	org := &client.Organization{
		Name: data.Name.ValueString(),
	}

	if !data.Description.IsNull() {
		org.Description = data.Description.ValueString()
	}
	if !data.MaxHosts.IsNull() {
		org.MaxHosts = int(data.MaxHosts.ValueInt64())
	}
	if !data.CustomVirtualEnv.IsNull() {
		org.CustomVirtualEnv = data.CustomVirtualEnv.ValueString()
	}

	createdOrg, err := r.client.CreateOrganization(org)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create organization, got error: %s", err),
		)
		return
	}

	data.ID = types.StringValue(strconv.Itoa(createdOrg.ID))
	// Re-hydrate other fields in case backend modified them
	data.Name = types.StringValue(createdOrg.Name)
	data.Description = types.StringValue(createdOrg.Description)
	data.MaxHosts = types.Int64Value(int64(createdOrg.MaxHosts))
	data.CustomVirtualEnv = types.StringValue(createdOrg.CustomVirtualEnv)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OrganizationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data OrganizationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID conversion", err.Error())
		return
	}

	org, err := r.client.GetOrganization(id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to read organization, got error: %s", err),
		)
		return
	}

	data.Name = types.StringValue(org.Name)
	data.Description = types.StringValue(org.Description)
	data.MaxHosts = types.Int64Value(int64(org.MaxHosts))
	data.CustomVirtualEnv = types.StringValue(org.CustomVirtualEnv)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OrganizationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data OrganizationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID conversion", err.Error())
		return
	}

	org := &client.Organization{
		ID:   id,
		Name: data.Name.ValueString(),
	}

	if !data.Description.IsNull() {
		org.Description = data.Description.ValueString()
	}
	if !data.MaxHosts.IsNull() {
		org.MaxHosts = int(data.MaxHosts.ValueInt64())
	}
	if !data.CustomVirtualEnv.IsNull() {
		org.CustomVirtualEnv = data.CustomVirtualEnv.ValueString()
	}

	updatedOrg, err := r.client.UpdateOrganization(org)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update organization, got error: %s", err),
		)
		return
	}

	data.Name = types.StringValue(updatedOrg.Name)
	data.Description = types.StringValue(updatedOrg.Description)
	data.MaxHosts = types.Int64Value(int64(updatedOrg.MaxHosts))
	data.CustomVirtualEnv = types.StringValue(updatedOrg.CustomVirtualEnv)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OrganizationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data OrganizationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID conversion", err.Error())
		return
	}

	err = r.client.DeleteOrganization(id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete organization, got error: %s", err),
		)
		return
	}
}

func (r *OrganizationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
