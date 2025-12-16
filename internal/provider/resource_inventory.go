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

var _ resource.Resource = &InventoryResource{}
var _ resource.ResourceWithImportState = &InventoryResource{}

func NewInventoryResource() resource.Resource {
	return &InventoryResource{}
}

type InventoryResource struct {
	client *client.Client
}

type InventoryResourceModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	OrganizationID types.String `tfsdk:"organization_id"`
	Kind           types.String `tfsdk:"kind"`
	HostFilter     types.String `tfsdk:"host_filter"`
	Variables      types.String `tfsdk:"variables"`
}

func (r *InventoryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_inventory"
}

func (r *InventoryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Inventory resource for AAP.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Inventory ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the inventory.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Description of the inventory.",
			},
			"organization_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "ID of the organization containing this inventory.",
			},
			"kind": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Kind of inventory. Empty for standard, 'smart' for smart inventory.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"host_filter": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Host filter for smart inventories.",
			},
			"variables": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Inventory variables in JSON or YAML format.",
			},
		},
	}
}

func (r *InventoryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", fmt.Sprintf("Expected *client.Client, got: %T.", req.ProviderData))
		return
	}

	r.client = client
}

func (r *InventoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data InventoryResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgID, _ := strconv.Atoi(data.OrganizationID.ValueString())

	inv := &client.Inventory{
		Name:         data.Name.ValueString(),
		Organization: orgID,
	}

	if !data.Description.IsNull() {
		inv.Description = data.Description.ValueString()
	}
	if !data.Kind.IsNull() {
		inv.Kind = data.Kind.ValueString()
	}
	if !data.HostFilter.IsNull() {
		inv.HostFilter = data.HostFilter.ValueString()
	}
	if !data.Variables.IsNull() {
		inv.Variables = data.Variables.ValueString()
	}

	createdInv, err := r.client.CreateInventory(inv)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create inventory: %s", err))
		return
	}

	data.ID = types.StringValue(strconv.Itoa(createdInv.ID))
	data.Name = types.StringValue(createdInv.Name)
	data.Description = types.StringValue(createdInv.Description)
	data.OrganizationID = types.StringValue(strconv.Itoa(createdInv.Organization))
	data.Kind = types.StringValue(createdInv.Kind)
	data.HostFilter = types.StringValue(createdInv.HostFilter)
	data.Variables = types.StringValue(createdInv.Variables)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InventoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data InventoryResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", err.Error())
		return
	}

	inv, err := r.client.GetInventory(id)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read inventory: %s", err))
		return
	}

	data.Name = types.StringValue(inv.Name)
	data.Description = types.StringValue(inv.Description)
	data.OrganizationID = types.StringValue(strconv.Itoa(inv.Organization))
	data.Kind = types.StringValue(inv.Kind)
	data.HostFilter = types.StringValue(inv.HostFilter)
	data.Variables = types.StringValue(inv.Variables)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InventoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data InventoryResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, _ := strconv.Atoi(data.ID.ValueString())
	orgID, _ := strconv.Atoi(data.OrganizationID.ValueString())

	inv := &client.Inventory{
		ID:           id,
		Name:         data.Name.ValueString(),
		Organization: orgID,
	}

	if !data.Description.IsNull() {
		inv.Description = data.Description.ValueString()
	}
	if !data.Kind.IsNull() {
		inv.Kind = data.Kind.ValueString()
	}
	if !data.HostFilter.IsNull() {
		inv.HostFilter = data.HostFilter.ValueString()
	}
	if !data.Variables.IsNull() {
		inv.Variables = data.Variables.ValueString()
	}

	updatedInv, err := r.client.UpdateInventory(inv)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update inventory: %s", err))
		return
	}

	data.Name = types.StringValue(updatedInv.Name)
	data.Description = types.StringValue(updatedInv.Description)
	data.OrganizationID = types.StringValue(strconv.Itoa(updatedInv.Organization))
	data.Kind = types.StringValue(updatedInv.Kind)
	data.HostFilter = types.StringValue(updatedInv.HostFilter)
	data.Variables = types.StringValue(updatedInv.Variables)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InventoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data InventoryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, _ := strconv.Atoi(data.ID.ValueString())
	if err := r.client.DeleteInventory(id); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete inventory: %s", err))
	}
}

func (r *InventoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
