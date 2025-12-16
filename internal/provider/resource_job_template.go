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

var _ resource.Resource = &JobTemplateResource{}
var _ resource.ResourceWithImportState = &JobTemplateResource{}

func NewJobTemplateResource() resource.Resource {
	return &JobTemplateResource{}
}

type JobTemplateResource struct {
	client *client.Client
}

type JobTemplateResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	JobType     types.String `tfsdk:"job_type"`
	InventoryID types.String `tfsdk:"inventory_id"`
	ProjectID   types.String `tfsdk:"project_id"`
	Playbook    types.String `tfsdk:"playbook"`
	Forks       types.Int64  `tfsdk:"forks"`
	Limit       types.String `tfsdk:"limit"`
	Verbosity   types.Int64  `tfsdk:"verbosity"`
	ExtraVars   types.String `tfsdk:"extra_vars"`
}

func (r *JobTemplateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_job_template"
}

func (r *JobTemplateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Job Template resource for AAP.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Job Template ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the job template.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Description of the job template.",
			},
			"job_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Type of job: 'run' or 'check'.",
			},
			"inventory_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "ID of the inventory to use.",
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "ID of the project to use.",
			},
			"playbook": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Playbook name to run.",
			},
			"forks": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Number of forks.",
			},
			"limit": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Limit to specific hosts.",
			},
			"verbosity": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Verbosity level (0-5).",
			},
			"extra_vars": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Extra variables in JSON/YAML format.",
			},
		},
	}
}

func (r *JobTemplateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *JobTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data JobTemplateResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	invID, _ := strconv.Atoi(data.InventoryID.ValueString())
	projID, _ := strconv.Atoi(data.ProjectID.ValueString())

	jt := &client.JobTemplate{
		Name:      data.Name.ValueString(),
		JobType:   data.JobType.ValueString(),
		Inventory: invID,
		Project:   projID,
		Playbook:  data.Playbook.ValueString(),
	}

	if !data.Description.IsNull() {
		jt.Description = data.Description.ValueString()
	}
	if !data.Forks.IsNull() {
		jt.Forks = int(data.Forks.ValueInt64())
	}
	if !data.Limit.IsNull() {
		jt.Limit = data.Limit.ValueString()
	}
	if !data.Verbosity.IsNull() {
		jt.Verbosity = int(data.Verbosity.ValueInt64())
	}
	if !data.ExtraVars.IsNull() {
		jt.ExtraVars = data.ExtraVars.ValueString()
	}

	createdJt, err := r.client.CreateJobTemplate(jt)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create job template: %s", err))
		return
	}

	data.ID = types.StringValue(strconv.Itoa(createdJt.ID))
	data.Name = types.StringValue(createdJt.Name)
	data.JobType = types.StringValue(createdJt.JobType)
	data.InventoryID = types.StringValue(strconv.Itoa(createdJt.Inventory))
	data.ProjectID = types.StringValue(strconv.Itoa(createdJt.Project))
	data.Playbook = types.StringValue(createdJt.Playbook)
	
	if createdJt.Description != "" {
		data.Description = types.StringValue(createdJt.Description)
	}
	data.Forks = types.Int64Value(int64(createdJt.Forks))
	if createdJt.Limit != "" {
		data.Limit = types.StringValue(createdJt.Limit)
	}
	data.Verbosity = types.Int64Value(int64(createdJt.Verbosity))
	if createdJt.ExtraVars != "" {
		data.ExtraVars = types.StringValue(createdJt.ExtraVars)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *JobTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data JobTemplateResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, _ := strconv.Atoi(data.ID.ValueString())

	jt, err := r.client.GetJobTemplate(id)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read job template: %s", err))
		return
	}

	data.Name = types.StringValue(jt.Name)
	data.Description = types.StringValue(jt.Description)
	data.JobType = types.StringValue(jt.JobType)
	data.InventoryID = types.StringValue(strconv.Itoa(jt.Inventory))
	data.ProjectID = types.StringValue(strconv.Itoa(jt.Project))
	data.Playbook = types.StringValue(jt.Playbook)
	data.Forks = types.Int64Value(int64(jt.Forks))
	data.Limit = types.StringValue(jt.Limit)
	data.Verbosity = types.Int64Value(int64(jt.Verbosity))
	data.ExtraVars = types.StringValue(jt.ExtraVars)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *JobTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data JobTemplateResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, _ := strconv.Atoi(data.ID.ValueString())
	invID, _ := strconv.Atoi(data.InventoryID.ValueString())
	projID, _ := strconv.Atoi(data.ProjectID.ValueString())

	jt := &client.JobTemplate{
		ID:        id,
		Name:      data.Name.ValueString(),
		JobType:   data.JobType.ValueString(),
		Inventory: invID,
		Project:   projID,
		Playbook:  data.Playbook.ValueString(),
	}

	if !data.Description.IsNull() {
		jt.Description = data.Description.ValueString()
	}
	if !data.Forks.IsNull() {
		jt.Forks = int(data.Forks.ValueInt64())
	}
	if !data.Limit.IsNull() {
		jt.Limit = data.Limit.ValueString()
	}
	if !data.Verbosity.IsNull() {
		jt.Verbosity = int(data.Verbosity.ValueInt64())
	}
	if !data.ExtraVars.IsNull() {
		jt.ExtraVars = data.ExtraVars.ValueString()
	}

	updatedJt, err := r.client.UpdateJobTemplate(jt)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update job template: %s", err))
		return
	}

	data.Name = types.StringValue(updatedJt.Name)
	data.Description = types.StringValue(updatedJt.Description)
	data.JobType = types.StringValue(updatedJt.JobType)
	data.InventoryID = types.StringValue(strconv.Itoa(updatedJt.Inventory))
	data.ProjectID = types.StringValue(strconv.Itoa(updatedJt.Project))
	data.Playbook = types.StringValue(updatedJt.Playbook)
	data.Forks = types.Int64Value(int64(updatedJt.Forks))
	data.Limit = types.StringValue(updatedJt.Limit)
	data.Verbosity = types.Int64Value(int64(updatedJt.Verbosity))
	data.ExtraVars = types.StringValue(updatedJt.ExtraVars)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *JobTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data JobTemplateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, _ := strconv.Atoi(data.ID.ValueString())
	if err := r.client.DeleteJobTemplate(id); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete job template: %s", err))
	}
}

func (r *JobTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
