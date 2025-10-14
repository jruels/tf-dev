package provider

import (
	"context"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type FileResource struct{}

var _ resource.Resource = &FileResource{}

func NewFileResource() resource.Resource {
	return &FileResource{}
}

func (r *FileResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "filemanager_file"
}

func (r *FileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a file and flags large ones.",
		Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{
				Required:    true,
				Description: "Full path to the file.",
			},
			"is_large": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the file is larger than 1GB.",
			},
		},
	}
}

func (r *FileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan struct {
		Path    types.String `tfsdk:"path"`
		IsLarge types.Bool   `tfsdk:"is_large"`
	}
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	info, err := os.Stat(plan.Path.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("File Error", err.Error())
		return
	}

	if info.IsDir() {
		resp.Diagnostics.AddError("Invalid File Path", "The specified path is a directory, not a file.")
		return
	}

	log.Println("************************")
	isLarge := true //info.Size() > 1<<30
	if isLarge {
		log.Printf("[INFO] Large file detected: %s (%d bytes)", plan.Path.ValueString(), info.Size())
	}

	plan.IsLarge = types.BoolValue(isLarge)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *FileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state struct {
		Path    types.String `tfsdk:"path"`
		IsLarge types.Bool   `tfsdk:"is_large"`
	}
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	info, err := os.Stat(state.Path.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("File Read Error", err.Error())
		return
	}

	if info.IsDir() {
		resp.Diagnostics.AddError("Invalid File Path", "The specified path is a directory, not a file.")
		return
	}

	isLarge := info.Size() > 1<<30
	state.IsLarge = types.BoolValue(isLarge)

	if isLarge {
		log.Printf("[INFO] Large file still exists: %s (%d bytes)", state.Path.ValueString(), info.Size())
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *FileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan struct {
		Path    types.String `tfsdk:"path"`
		IsLarge types.Bool   `tfsdk:"is_large"`
	}
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	info, err := os.Stat(plan.Path.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("File Error", err.Error())
		return
	}

	if info.IsDir() {
		resp.Diagnostics.AddError("Invalid File Path", "The specified path is a directory, not a file.")
		return
	}

	isLarge := info.Size() > 1<<30
	if isLarge {
		log.Printf("[INFO] Large file (update): %s (%d bytes)", plan.Path.ValueString(), info.Size())
	}

	plan.IsLarge = types.BoolValue(isLarge)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *FileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
