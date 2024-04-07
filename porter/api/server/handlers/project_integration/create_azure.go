package project_integration

import (
	"net/http"

	"github.com/porter-dev/porter/internal/telemetry"

	"connectrpc.com/connect"

	porterv1 "github.com/porter-dev/api-contracts/generated/go/porter/v1"

	"github.com/porter-dev/porter/api/server/handlers"
	"github.com/porter-dev/porter/api/server/shared"
	"github.com/porter-dev/porter/api/server/shared/apierrors"
	"github.com/porter-dev/porter/api/server/shared/config"
	"github.com/porter-dev/porter/api/types"
	"github.com/porter-dev/porter/internal/models"
	ints "github.com/porter-dev/porter/internal/models/integrations"
)

type CreateAzureHandler struct {
	handlers.PorterHandlerReadWriter
}

func NewCreateAzureHandler(
	config *config.Config,
	decoderValidator shared.RequestDecoderValidator,
	writer shared.ResultWriter,
) *CreateAzureHandler {
	return &CreateAzureHandler{
		PorterHandlerReadWriter: handlers.NewDefaultPorterHandler(config, decoderValidator, writer),
	}
}

func (p *CreateAzureHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, span := telemetry.NewSpan(r.Context(), "serve-create-azure-connection")
	defer span.End()

	user, _ := ctx.Value(types.UserScope).(*models.User)
	project, _ := ctx.Value(types.ProjectScope).(*models.Project)

	telemetry.WithAttributes(span, telemetry.AttributeKV{Key: "project-id", Value: project.ID})

	request := &types.CreateAzureRequest{}

	if ok := p.DecodeAndValidate(w, r, request); !ok {
		err := telemetry.Error(ctx, span, nil, "error decoding and validating request")
		p.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusBadRequest))
		return
	}

	az := CreateAzureIntegration(request, project.ID, user.ID)

	az, err := p.Repo().AzureIntegration().CreateAzureIntegration(az)
	if err != nil {
		err = telemetry.Error(ctx, span, err, "error creating azure integration")
		p.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusInternalServerError))
		return
	}

	res := types.CreateAzureResponse{
		AzureIntegration: az.ToAzureIntegrationType(),
	}

	if p.Config().ClusterControlPlaneClient == nil {
		err := telemetry.Error(ctx, span, nil, "cluster control plane client cannot be nil")
		p.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusInternalServerError))
		return
	}

	req := connect.NewRequest(&porterv1.SaveAzureCredentialsRequest{
		ProjectId:              int64(project.ID),
		ClientId:               request.AzureClientID,
		SubscriptionId:         request.AzureSubscriptionID,
		TenantId:               request.AzureTenantID,
		ServicePrincipalSecret: []byte(request.ServicePrincipalKey),
	})
	resp, err := p.Config().ClusterControlPlaneClient.SaveAzureCredentials(ctx, req)
	if err != nil {
		err = telemetry.Error(ctx, span, err, "error saving azure credentials")
		p.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusInternalServerError))
		return
	}

	if resp.Msg == nil {
		err = telemetry.Error(ctx, span, nil, "SaveAzureCredentials response message is nil")
		p.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusInternalServerError))
		return
	}
	if resp.Msg.CredentialsIdentifier == "" {
		err = telemetry.Error(ctx, span, nil, "SaveAzureCredentials response credentials identifier is empty")
		p.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusInternalServerError))
		return
	}

	res.CloudProviderCredentialIdentifier = resp.Msg.CredentialsIdentifier

	p.WriteResult(w, r, res)
}

func CreateAzureIntegration(request *types.CreateAzureRequest, projectID, userID uint) *ints.AzureIntegration {
	resp := &ints.AzureIntegration{
		UserID:                 userID,
		ProjectID:              projectID,
		AzureClientID:          request.AzureClientID,
		AzureSubscriptionID:    request.AzureSubscriptionID,
		AzureTenantID:          request.AzureTenantID,
		ServicePrincipalSecret: []byte(request.ServicePrincipalKey),
	}

	return resp
}
