var config = {
	"borrowPeriodDays": 15,
	"numBookOneCanBorrow": 3
};

var hoc5App = angular.module("hoc5App", ["ngRoute", "ui.bootstrap", "ui.utils"]);

hoc5App.config(["$interpolateProvider", function($interpolateProvider) {
	$interpolateProvider.startSymbol('[[');
	$interpolateProvider.endSymbol(']]');
}]);

hoc5App.run(function($rootScope, $templateCache) {
    $rootScope.$on('$routeChangeStart', function(event, next, current) {
        if (typeof(current) !== 'undefined'){
            $templateCache.remove(current.templateUrl);
        }
    });
});

hoc5App.config(["$routeProvider", function($routeProvider) {
	$routeProvider.
		when("/menu", {
			templateUrl: "partials/menu.html",
			controller: "MenuCtrl"
		}).
		when("/book/borrow", {
			templateUrl: "partials/book/borrow.html",
			controller: "BookBorrowCtrl"
		}).
		when("/book/return", {
			templateUrl: "partials/book/return.html",
			controller: "BookReturnCtrl"
		}).
		when("/book/manage", {
			templateUrl: "partials/book/manage.html",
			controller: "BookManageCtrl"
		}).
		when("/book/new", {
			templateUrl: "partials/book/edit.html",
			controller: "BookNewCtrl"
		}).
		when("/book/:barcode/edit", {
			templateUrl: "partials/book/edit.html",
			controller: "BookEditCtrl"
		}).
		when("/book/overdue", {
			templateUrl: "partials/book/overdue.html",
			controller: "BookOverdueCtrl"
		}).
		when("/borrower/manage", {
			templateUrl: "partials/borrower/manage.html",
			controller: "BorrowerManageCtrl"
		}).
		when("/borrower/new", {
			templateUrl: "partials/borrower/edit.html",
			controller: "BorrowerNewCtrl"
		}).
		when("/borrower/:phone/edit", {
			templateUrl: "partials/borrower/edit.html",
			controller: "BorrowerEditCtrl"
		}).
		when("/book/stock", {
			templateUrl: "partials/book/stock.html",
			controller: "BookStockCtrl"
		}).
		when("/book/stock/:begin/:end", {
			templateUrl: "partials/book/stock-entry.html",
			controller: "BookStockEntryCtrl"
		}).
		when("/gen-barcode", {
			templateUrl: "partials/gen-barcode.html",
			controller: "GenBarcodeCtrl"
		}).
		otherwise({
			redirectTo: "/menu"
		});
}]);

hoc5App.controller('BookReturnCtrl', [
	'$scope', '$http', '$location', function($scope, $http, $location){
	$scope.$parent.page = {
		title: "Return Books"
	};
	$scope.suggest = {
		book: false
	};
	$scope.ShowSuggest = function(show){
		$scope.suggest.book = show;
	};
	$scope.return_ = {
		barcode: ""
	};
	$scope.borrower = null;
	$scope.SetBarcode = function(barcode){
		$scope.return_.barcode = barcode;
		if (barcode) {
			$http({
				method: "GET",
				url: "/api/book/" + barcode + "/borrower"
			}).success(function(data, status){
				$scope.borrower = data.Borrower;
			}).error(function(data, status){
				$scope.$parent.errors.push(data);
				$scope.borrower = null;
			});
		} else {
			$scope.borrower = null;
		}
	};
	$scope.IsBarcodeValid = function() {
		return !!$scope.return_.barcode;
	};
	$scope.Submit = function(){
		if ($scope.inProgress) {
			return;
		}
		$scope.$parent.page.errors = [];
		if (!$scope.IsBarcodeValid()) {
			$scope.$parent.page.errors.push('Barcode is invalid');
		}
		if ($scope.$parent.page.errors.length) {
			return;
		}
		$scope.inProgress = true;
		$http({
			method: "POST",
			url: "/api/book/return",
			headers: {'Content-Type': 'application/x-www-form-urlencoded'},
			data: $.param($scope.return_)
		}).success(function(data, status){
			// TODO: This seems not working. Fix it later
			$scope.$parent.page.messages = ["Return book successful"];
			$location.path("/menu");
		}).error(function(data, status){
			$scope.$parent.page.errors = [data];
			$scope.inProgress = false;
		});
	};
}]);

hoc5App.controller('BookSearchCtrl', ["$scope", "$http", function($scope, $http){
	$scope.books = [];
	$scope.query = "";
	$scope.error = "";
	$scope.currentPage = 1;
	$scope.numResults = 0;
	$scope.itemsPerPage = 10;
	$scope.SetQuery = function(query) {
		$scope.query = query;
	};
	$scope.Search = function(page) {
		var req;
		if ($scope.query) {
			req = $http({
				method: "GET",
				url: "/api/book/search",
				params: {"q": $scope.query, "page": page}
			});
		} else {
			req = $http({
				method: "GET",
				url: "/api/book/list",
				params: {"page": page}
			});
		}
		req.success(function(data, status){
			$scope.books = data.Books;
			$scope.currentPage = data.Page;
			$scope.numResults = data.NumResults;
			$scope.error = "";
		}).error(function(data, status){
			$scope.books = [];
			$scope.currentPage = 1;
			$scope.numResults = 0;
			$scope.error = data;
		});
	};
}]);

hoc5App.controller('BorrowerSearchCtrl', ["$scope", "$http", function($scope, $http){
	$scope.borrowers = [];
	$scope.query = "";
	$scope.error = "";
	$scope.currentPage = 1;
	$scope.numResults = 0;
	$scope.itemsPerPage = 5;
	$scope.SetQuery = function(query) {
		$scope.query = query;
	};
	$scope.CanBorrow = function(borrower) {
		return borrower.NumBorrowed < config.numBookOneCanBorrow;
	};
	$scope.Search = function(page) {
		var req;
		if ($scope.query) {
			req = $http({
				method: "GET",
				url: "/api/borrower/search",
				params: {"q": $scope.query, "page": page}
			});
		} else {
			req = $http({
				method: "GET",
				url: "/api/borrower/list",
				params: {"page": page}
			});
		}
		req.success(function(data, status){
			$scope.borrowers = data.Borrowers;
			$scope.currentPage = data.Page;
			$scope.numResults = data.NumResults;
			$scope.error = "";
		}).error(function(data, status){
			$scope.borrowers = [];
			$scope.currentPage = 1;
			$scope.numResults = 0;
			$scope.error = data;
		});
	};
}]);




