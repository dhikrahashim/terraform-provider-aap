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

var _ resource.Resource = &CredentialScmResource{}
var _ resource.ResourceWithImportState = &CredentialScmResource{}

func NewCredentialScmResource() resource.Resource {
	return &CredentialScmResource{}
}

type CredentialScmResource struct {
	client *client.Client
}

type CredentialScmResourceModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	OrganizationID types.String `tfsdk:"organization_id"`
	Username       types.String `tfsdk:"username"`
	Password       types.String `tfsdk:"password"`
	SSHKeyData     types.String `tfsdk:"ssh_key_data"`
	SSHKeyUnlock   types.String `tfsdk:"ssh_key_unlock"`
}

func (r *CredentialScmResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential_scm"
}

func (r *CredentialScmResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Source Control credential for accessing SCM repositories.",
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
			"organization_id": schema.StringAttribute{
				Required: true,
			},
			"username": schema.StringAttribute{
				Optional: true,
			},
			"password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"ssh_key_data": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"ssh_key_unlock": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func (r *CredentialScmResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CredentialScmResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CredentialScmResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgID, _ := strconv.Atoi(data.OrganizationID.ValueString())

	cred := &client.Credential{
		Name:           data.Name.ValueString(),
		Description:    data.Description.ValueString(),
		Organization:   orgID,
		CredentialType: 2, // SCM credential type
		Inputs: client.CredentialInputs{
			Username:     data.Username.ValueString(),
			Password:     data.Password.ValueString(),
			SSHKeyData:   data.SSHKeyData.ValueString(),
			SSHKeyUnlock: data.SSHKeyUnlock.ValueString(),
		},
	}

	created, err := r.client.CreateCredential(cred)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create SCM credential: %s", err))
		return
	}

	data.ID = types.StringValue(strconv.Itoa(created.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CredentialScmResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CredentialScmResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, _ := strconv.Atoi(data.ID.ValueString())
	cred, err := r.client.GetCredential(id)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read SCM credential: %s", err))
		return
	}

	data.Name = types.StringValue(cred.Name)
	data.Description = types.StringValue(cred.Description)
	data.OrganizationID = types.StringValue(strconv.Itoa(cred.Organization))
	data.Username = types.StringValue(cred.Inputs.Username)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CredentialScmResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data CredentialScmResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, _ := strconv.Atoi(data.ID.ValueString())
	orgID, _ := strconv.Atoi(data.OrganizationID.ValueString())

	cred := &client.Credential{
		ID:             id,
		Name:           data.Name.ValueString(),
		Description:    data.Description.ValueString(),
		Organization:   orgID,
		CredentialType: 2,
		Inputs: client.CredentialInputs{
			Username:     data.Username.ValueString(),
			Password:     data.Password.ValueString(),
			SSHKeyData:   data.SSHKeyData.ValueString(),
			SSHKeyUnlock: data.SSHKeyUnlock.ValueString(),
		},
	}

	_, err := r.client.UpdateCredential(cred)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update SCM credential: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CredentialScmResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CredentialScmResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, _ := strconv.Atoi(data.ID.ValueString())
	if err := r.client.DeleteCredential(id); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete SCM credential: %s", err))
	}
}

func (r *CredentialScmResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
