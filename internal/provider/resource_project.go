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

var _ resource.Resource = &ProjectResource{}
var _ resource.ResourceWithImportState = &ProjectResource{}

func NewProjectResource() resource.Resource {
	return &ProjectResource{}
}

type ProjectResource struct {
	client *client.Client
}

type ProjectResourceModel struct {
	ID                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	Description           types.String `tfsdk:"description"`
	OrganizationID        types.String `tfsdk:"organization_id"`
	ScmType               types.String `tfsdk:"scm_type"`
	ScmUrl                types.String `tfsdk:"scm_url"`
	ScmBranch             types.String `tfsdk:"scm_branch"`
	ScmCredentialID       types.String `tfsdk:"scm_credential_id"`
	ScmClean              types.Bool   `tfsdk:"scm_clean"`
	ScmDeleteOnUpdate     types.Bool   `tfsdk:"scm_delete_on_update"`
	ScmUpdateOnLaunch     types.Bool   `tfsdk:"scm_update_on_launch"`
	ScmUpdateCacheTimeout types.Int64  `tfsdk:"scm_update_cache_timeout"`
	LocalPath             types.String `tfsdk:"local_path"`
}

func (r *ProjectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (r *ProjectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Project resource for AAP. Projects represent SCM repositories containing playbooks.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the project.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Description of the project.",
			},
			"organization_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Organization ID.",
			},
			"scm_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "SCM type: '' (manual), 'git', 'hg', 'svn'.",
			},
			"scm_url": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "SCM repository URL.",
			},
			"scm_branch": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Branch, tag, or commit to checkout.",
			},
			"scm_credential_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "SCM credential ID.",
			},
			"scm_clean": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Clean the repository before syncing.",
			},
			"scm_delete_on_update": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Delete local modifications before updating.",
			},
			"scm_update_on_launch": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Update the project when a job is launched.",
			},
			"scm_update_cache_timeout": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Cache timeout for SCM updates.",
			},
			"local_path": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Local path for manual projects.",
			},
		},
	}
}

func (r *ProjectResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ProjectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgID, _ := strconv.Atoi(data.OrganizationID.ValueString())
	credID, _ := strconv.Atoi(data.ScmCredentialID.ValueString())

	p := &client.Project{
		Name:         data.Name.ValueString(),
		Organization: orgID,
		ScmType:      data.ScmType.ValueString(),
	}
	if !data.Description.IsNull() {
		p.Description = data.Description.ValueString()
	}
	if !data.ScmUrl.IsNull() {
		p.ScmUrl = data.ScmUrl.ValueString()
	}
	if !data.ScmBranch.IsNull() {
		p.ScmBranch = data.ScmBranch.ValueString()
	}
	if credID > 0 {
		p.ScmCredential = credID
	}
	if !data.ScmClean.IsNull() {
		p.ScmClean = data.ScmClean.ValueBool()
	}
	if !data.ScmDeleteOnUpdate.IsNull() {
		p.ScmDeleteOnUpdate = data.ScmDeleteOnUpdate.ValueBool()
	}
	if !data.ScmUpdateOnLaunch.IsNull() {
		p.ScmUpdateOnLaunch = data.ScmUpdateOnLaunch.ValueBool()
	}
	if !data.ScmUpdateCacheTimeout.IsNull() {
		p.ScmUpdateCacheTimeout = int(data.ScmUpdateCacheTimeout.ValueInt64())
	}
	if !data.LocalPath.IsNull() {
		p.LocalPath = data.LocalPath.ValueString()
	}

	created, err := r.client.CreateProject(p)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create project: %s", err))
		return
	}

	data.ID = types.StringValue(strconv.Itoa(created.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ProjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, _ := strconv.Atoi(data.ID.ValueString())
	p, err := r.client.GetProject(id)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read project: %s", err))
		return
	}

	data.Name = types.StringValue(p.Name)
	data.Description = types.StringValue(p.Description)
	data.OrganizationID = types.StringValue(strconv.Itoa(p.Organization))
	data.ScmType = types.StringValue(p.ScmType)
	data.ScmUrl = types.StringValue(p.ScmUrl)
	data.ScmBranch = types.StringValue(p.ScmBranch)
	if p.ScmCredential > 0 {
		data.ScmCredentialID = types.StringValue(strconv.Itoa(p.ScmCredential))
	}
	data.ScmClean = types.BoolValue(p.ScmClean)
	data.ScmDeleteOnUpdate = types.BoolValue(p.ScmDeleteOnUpdate)
	data.ScmUpdateOnLaunch = types.BoolValue(p.ScmUpdateOnLaunch)
	data.ScmUpdateCacheTimeout = types.Int64Value(int64(p.ScmUpdateCacheTimeout))
	data.LocalPath = types.StringValue(p.LocalPath)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ProjectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, _ := strconv.Atoi(data.ID.ValueString())
	orgID, _ := strconv.Atoi(data.OrganizationID.ValueString())
	credID, _ := strconv.Atoi(data.ScmCredentialID.ValueString())

	p := &client.Project{
		ID:           id,
		Name:         data.Name.ValueString(),
		Organization: orgID,
		ScmType:      data.ScmType.ValueString(),
	}
	if !data.Description.IsNull() {
		p.Description = data.Description.ValueString()
	}
	if !data.ScmUrl.IsNull() {
		p.ScmUrl = data.ScmUrl.ValueString()
	}
	if !data.ScmBranch.IsNull() {
		p.ScmBranch = data.ScmBranch.ValueString()
	}
	if credID > 0 {
		p.ScmCredential = credID
	}
	if !data.ScmClean.IsNull() {
		p.ScmClean = data.ScmClean.ValueBool()
	}
	if !data.ScmDeleteOnUpdate.IsNull() {
		p.ScmDeleteOnUpdate = data.ScmDeleteOnUpdate.ValueBool()
	}
	if !data.ScmUpdateOnLaunch.IsNull() {
		p.ScmUpdateOnLaunch = data.ScmUpdateOnLaunch.ValueBool()
	}
	if !data.ScmUpdateCacheTimeout.IsNull() {
		p.ScmUpdateCacheTimeout = int(data.ScmUpdateCacheTimeout.ValueInt64())
	}

	_, err := r.client.UpdateProject(p)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update project: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ProjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, _ := strconv.Atoi(data.ID.ValueString())
	if err := r.client.DeleteProject(id); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete project: %s", err))
	}
}

func (r *ProjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
