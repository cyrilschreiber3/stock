package router

import (
	"github.com/cyrilschreiber3/stock/handlers"
	"github.com/cyrilschreiber3/stock/routes"
	"github.com/gin-gonic/gin"
)

func registerCoreRoutes(r *gin.Engine) {
	routes.Index.Register(r, handlers.Index())
}

func registerApiRoutes(r *gin.Engine) {
	routes.ApiHealth.Register(r, handlers.HandleGetHealth())
}

func registerProductRoutes(r *gin.Engine) {
	routes.ProductList.Register(r, handlers.HandleGetProducts())
	routes.ProductOptions.Register(r, handlers.HandleGetProductOptions())
	routes.ProductSearch.Register(r, handlers.HandleSearchProducts())
	routes.ProductSearchTransactions.Register(r, handlers.HandleSearchTransactionsForProduct())
	routes.ProductSearchInventoryLots.Register(r, handlers.HandleSearchInventoryLotsForProduct())
	routes.ProductCreateForm.Register(r, handlers.HandleShowCreateProductForm())
	routes.ProductCreate.Register(r, handlers.HandleCreateProduct())
	routes.ProductDetails.Register(r, handlers.HandleGetProductDetails())
	routes.ProductField.Register(r, handlers.HandleGetProductFieldValue())
	routes.ProductEditForm.Register(r, handlers.HandleShowUpdateProductForm())
	routes.ProductUpdate.Register(r, handlers.HandleUpdateProduct())
	routes.ProductDelete.Register(r, handlers.HandleDeleteProduct())
}

func registerBrandRoutes(r *gin.Engine) {
	routes.BrandOptions.Register(r, handlers.HandleGetProductBrandOptions())
	routes.BrandDetails.Register(r, handlers.HandleGetBrandDetails())
	routes.BrandSearchProducts.Register(r, handlers.HandleSearchProductsByBrand())
}

func registerCategoryRoutes(r *gin.Engine) {
	routes.CategoryList.Register(r, handlers.HandleGetCategories())
	routes.CategoryOptions.Register(r, handlers.HandleGetCategoryOptions())
	routes.CategorySearch.Register(r, handlers.HandleSearchCategories())
	routes.CategoryCreateForm.Register(r, handlers.HandleShowCreateCategoryForm())
	routes.CategoryCreate.Register(r, handlers.HandleCreateCategory())
	routes.CategoryDetails.Register(r, handlers.HandleGetCategoryDetails())
	routes.CategorySearchProducts.Register(r, handlers.HandleSearchProductsByCategory())
	routes.CategoryEditForm.Register(r, handlers.HandleShowUpdateCategoryForm())
	routes.CategoryUpdate.Register(r, handlers.HandleUpdateCategory())
	routes.CategoryDelete.Register(r, handlers.HandleDeleteCategory())
}

func registerSubcategoryRoutes(r *gin.Engine) {
	routes.SubcategoryOptions.Register(r, handlers.HandleGetSubcategoryOptions())
	routes.SubcategorySearch.Register(r, handlers.HandleSearchSubcategories())
	routes.SubcategoryCreateForm.Register(r, handlers.HandleShowCreateSubcategoryForm())
	routes.SubcategoryCreate.Register(r, handlers.HandleCreateSubcategory())
	routes.SubcategoryDetails.Register(r, handlers.HandleGetSubcategoryDetails())
	routes.SubcategorySearchProducts.Register(r, handlers.HandleSearchProductsBySubcategory())
	routes.SubcategoryEditForm.Register(r, handlers.HandleShowUpdateSubcategoryForm())
	routes.SubcategoryUpdate.Register(r, handlers.HandleUpdateSubcategory())
	routes.SubcategoryDelete.Register(r, handlers.HandleDeleteSubcategory())
}

func registerSupplierRoutes(r *gin.Engine) {
	routes.SupplierList.Register(r, handlers.HandleGetSuppliers())
	routes.SupplierOptions.Register(r, handlers.HandleGetSupplierOptions())
	routes.SupplierSearch.Register(r, handlers.HandleSearchSuppliers())
	routes.SupplierCreateForm.Register(r, handlers.HandleShowCreateSupplierForm())
	routes.SupplierCreate.Register(r, handlers.HandleCreateSupplier())
	routes.SupplierDetails.Register(r, handlers.HandleGetSupplierDetails())
	routes.SupplierSearchProducts.Register(r, handlers.HandleSearchProductsBySupplier())
	routes.SupplierEditForm.Register(r, handlers.HandleShowUpdateSupplierForm())
	routes.SupplierUpdate.Register(r, handlers.HandleUpdateSupplier())
	routes.SupplierDelete.Register(r, handlers.HandleDeleteSupplier())
}

func registerTransactionRoutes(r *gin.Engine) {
	routes.TransactionList.Register(r, handlers.HandleGetTransactions())
	routes.TransactionSearch.Register(r, handlers.HandleSearchTransactions())
	routes.TransactionCreateForm.Register(r, handlers.HandleShowCreateTransactionForm())
	routes.TransactionCreate.Register(r, handlers.HandleCreateTransaction())
	routes.TransactionSearchProductsForm.Register(r, handlers.HandleShowSearchProductsForTransactionItems())
	routes.TransactionSearchProducts.Register(r, handlers.HandleSearchProductsForTransactionItems())
	routes.TransactionSelectProduct.Register(r, handlers.HandleSelectProductForTransactionItem())
	routes.TransactionDetails.Register(r, handlers.HandleGetTransactionDetails())
	routes.TransactionEditForm.Register(r, handlers.HandleShowUpdateTransactionForm())
	routes.TransactionUpdate.Register(r, handlers.HandleUpdateTransaction())
	routes.TransactionDelete.Register(r, handlers.HandleDeleteTransaction())
	routes.TransactionApply.Register(r, handlers.HandleApplyTransaction())
}

func registerTransactionItemRoutes(r *gin.Engine) {
	routes.TransactionItemSearch.Register(r, handlers.HandleSearchTransactionItems())
	routes.TransactionItemCreate.Register(r, handlers.HandleCreateTransactionItem())
	routes.TransactionItemEditForm.Register(r, handlers.HandleShowUpdateTransactionItemForm())
	routes.TransactionItemUpdate.Register(r, handlers.HandleUpdateTransactionItem())
	routes.TransactionItemDelete.Register(r, handlers.HandleDeleteTransactionItem())
}

func registerInventoryRoutes(r *gin.Engine) {
	routes.InventoryList.Register(r, handlers.HandleGetInventory())
	// routes.InventoryDetails.Register(r, handlers.HandleGetInventoryDetails())
	routes.InventorySearch.Register(r, handlers.HandleSearchInventory())
}
