package template

import (
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"kubesphere.io/devops/pkg/api"
	"kubesphere.io/devops/pkg/api/devops/v1alpha1"
	"kubesphere.io/devops/pkg/constants"
	kapisv1alpha1 "kubesphere.io/devops/pkg/kapis/devops/v1alpha1/common"
	"net/http"
)

var (
	// DevopsPathParameter is a path parameter definition for devops.
	DevopsPathParameter = restful.PathParameter("devops", "DevOps project name")
	// TemplatePathParameter is a path parameter definition for template.
	TemplatePathParameter = restful.PathParameter("template", "Template name")
)

func Routers(service *restful.WebService, options *kapisv1alpha1.Options) {
	handler := newHandler(options)
	service.Route(service.GET("/devops/{devops}/templates").
		To(handler.handleQuery).
		Param(DevopsPathParameter).
		Doc("Query templates for a DevOps Project.").
		Returns(http.StatusOK, api.StatusOK, api.ListResult{Items: []interface{}{}}).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.DevOpsTemplateTag}))

	service.Route(service.POST("/devops/{devops}/templates").
		// TODO Add handler to this route
		To(handler.handleCreate).
		Param(DevopsPathParameter).
		Doc("Create a template for a DevOps Project.").
		Reads(v1alpha1.Template{}).
		Returns(http.StatusCreated, api.StatusOK, v1alpha1.Template{}).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.DevOpsTemplateTag}))

	service.Route(service.GET("/devops/{devops}/templates/{template}").
		To(handler.handleGet).
		Param(DevopsPathParameter).
		Param(TemplatePathParameter).
		Doc("Get template").
		Returns(http.StatusOK, api.StatusOK, v1alpha1.Template{}).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.DevOpsTemplateTag}))

	service.Route(service.POST("/devops/{devops}/templates/{template}/render").
		To(handler.handleRender).
		Param(DevopsPathParameter).
		Param(TemplatePathParameter).
		Doc("Render template using the given parameters and return render result into annotations inside template").
		Reads([]Parameter{}).
		Returns(http.StatusOK, api.StatusOK, v1alpha1.Template{}).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.DevOpsTemplateTag}))

	service.Route(service.PUT("/devops/{devops}/templates/{template}").
		// TODO Add handler to this route
		To(handler.handleQuery).
		Param(DevopsPathParameter).
		Param(TemplatePathParameter).
		Doc("Update template").
		Reads(v1alpha1.Template{}).
		Returns(http.StatusOK, api.StatusOK, v1alpha1.Template{}).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.DevOpsTemplateTag}))

	service.Route(service.DELETE("/devops/{devops}/templates/{template}").
		// TODO Add handler to this route
		To(handler.handleQuery).
		Param(DevopsPathParameter).
		Param(TemplatePathParameter).
		Doc("Delete a template").
		Returns(http.StatusOK, api.StatusOK, v1alpha1.Template{}).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.DevOpsTemplateTag}))
}
