package routes

import "github.com/cyrilschreiber3/stock/handlers"

var CategoryGroup = Group("/categories")

var CategoryList = GET[NoParams]("CategoryList", CategoryGroup, "", handlers.HandleGetCategories())

var CategoryOptions = GET[NoParams]("CategoryOptions", CategoryGroup, "/options", handlers.HandleGetCategoryOptions())

var CategoryCreateForm = GET[NoParams]("CategoryCreateForm", CategoryGroup, "/create", handlers.HandleShowCreateCategoryForm())

var CategoryCreate = POST[NoParams]("CategoryCreate", CategoryGroup, "/create", handlers.HandleCreateCategory())

var CategoryDetailsRouteParams = WithRouteParams[SimpleUUIDParam](ParamID)
var CategoryDetails = GET("CategoryDetails", CategoryGroup, "/:id", handlers.HandleGetCategoryDetails(), CategoryDetailsRouteParams)

var CategoryEditFormRouteParams = WithRouteParams[SimpleUUIDParam](ParamID)
var CategoryEditForm = GET("CategoryEditForm", CategoryGroup, "/:id/edit", handlers.HandleShowUpdateCategoryForm(), CategoryEditFormRouteParams)

var CategoryUpdateRouteParams = WithRouteParams[SimpleUUIDParam](ParamID)
var CategoryUpdate = PUT("CategoryUpdate", CategoryGroup, "/:id/update", handlers.HandleUpdateCategory(), CategoryUpdateRouteParams)

var CategoryDeleteRouteParams = WithRouteParams[SimpleUUIDParam](ParamID)
var CategoryDelete = DELETE("CategoryDelete", CategoryGroup, "/:id/delete", handlers.HandleDeleteCategory(), CategoryDeleteRouteParams)

var CategoryRoutes = ResourceRoutes{
	Group: CategoryGroup,
	Routes: []RouteHandler{
		CategoryList,
		CategoryOptions,
		CategoryCreateForm,
		CategoryCreate,
		CategoryDetails,
		CategoryEditForm,
		CategoryUpdate,
		CategoryDelete,
	},
}

var SubcategoryGroupDefaultParams = WithRouteParams[SimpleUUIDParam](ParamID)
var SubcategoryGroup = Group("/:id/subcategories", WithParent(CategoryGroup))

var SubcategoryList = GET("SubcategoryList", SubcategoryGroup, "", handlers.HandleGetSubcategories(), SubcategoryGroupDefaultParams)

var SubcategoryOptions = GET("SubcategoryOptions", SubcategoryGroup, "/options", handlers.HandleGetSubcategoryOptions(), SubcategoryGroupDefaultParams)

var SubcategoryCreateForm = GET("SubcategoryCreateForm", SubcategoryGroup, "/create", handlers.HandleShowCreateSubcategoryForm(), SubcategoryGroupDefaultParams)

var SubcategoryCreate = POST("SubcategoryCreate", SubcategoryGroup, "/create", handlers.HandleCreateSubcategory(), SubcategoryGroupDefaultParams)

// var SubcategoryDetailsRouteParams = WithRouteParams[SubcategoryUUIDParams](ParamID, ParamSubcategoryID)
// var SubcategoryDetails = GET("SubcategoryDetails", SubcategoryGroup, "/:subcategory_id", handlers.HandleGetSubcategoryDetails(), SubcategoryDetailsRouteParams)

var SubcategoryEditFormRouteParams = WithRouteParams[SubcategoryUUIDParams](ParamID, ParamSubcategoryID)
var SubcategoryEditForm = GET("SubcategoryEditForm", SubcategoryGroup, "/:subcategory_id/edit", handlers.HandleShowUpdateSubcategoryForm(), SubcategoryEditFormRouteParams)

var SubcategoryUpdateRouteParams = WithRouteParams[SubcategoryUUIDParams](ParamID, ParamSubcategoryID)
var SubcategoryUpdate = PUT("SubcategoryUpdate", SubcategoryGroup, "/:subcategory_id/update", handlers.HandleUpdateSubcategory(), SubcategoryUpdateRouteParams)

var SubcategoryDeleteRouteParams = WithRouteParams[SubcategoryUUIDParams](ParamID, ParamSubcategoryID)
var SubcategoryDelete = DELETE("SubcategoryDelete", SubcategoryGroup, "/:subcategory_id/delete", handlers.HandleDeleteSubcategory(), SubcategoryDeleteRouteParams)

var SubcategoryRoutes = ResourceRoutes{
	Group: SubcategoryGroup,
	Routes: []RouteHandler{
		SubcategoryList,
		SubcategoryOptions,
		SubcategoryCreateForm,
		SubcategoryCreate,
		// SubcategoryDetails,
		SubcategoryEditForm,
		SubcategoryUpdate,
		SubcategoryDelete,
	},
}
