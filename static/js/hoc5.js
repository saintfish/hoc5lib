var config = {
	"borrowPeriodDays": 15
};

function range(start, stop, step){
    if (typeof stop=='undefined') {
        // one param defined
        stop = start;
        start = 0;
    }
    if (typeof step=='undefined') {
        step = 1;
    }
    if ((step>0 && start>=stop) || (step<0 && start<=stop)) {
        return [];
    }
    var result = [];
    for (var i=start; step>0 ? i<stop : i>stop; i+=step) {
        result.push(i);
    }
    return result;
}

var hoc5App = angular.module("hoc5App", ["ngRoute"]);

hoc5App.config(["$interpolateProvider", function($interpolateProvider) {
	$interpolateProvider.startSymbol('[[');
	$interpolateProvider.endSymbol(']]');
}]);

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
		otherwise({
			redirectTo: "/menu"
		});
}]);

hoc5App.controller('MenuCtrl', ['$scope', function($scope){
	$scope.$parent.page = {
		title: "Menu"
	};
}]);

hoc5App.controller('BookBorrowCtrl', [
	'$scope', '$http', '$location', function($scope, $http, $location){
	$scope.$parent.page = {
		title: "Borrow Books"
	};
	$scope.returnDate = function(){
		var d = new Date();
		d.setDate(d.getDate() + config.borrowPeriodDays);
		return d;
	}();
	$scope.suggest = {
		book: false,
		borrower: false
	};
	$scope.borrow = {
		"barcode": "",
		"phone": "",
	};
	$scope.inProgress = false;

	$scope.IsBarcodeValid = function() {
		return !!$scope.borrow.barcode;
	};
	$scope.IsPhoneValid = function() {
		return !!$scope.borrow.phone;
	};
	$scope.SetBarcode = function(barcode) {
		$scope.borrow.barcode = barcode;
	};
	$scope.SetPhone = function(phone) {
		$scope.borrow.phone = phone;
	};
	$scope.ShowSuggest = function(b) {
		$scope.suggest.book = b;
		$scope.suggest.borrower = b;
	};
	$scope.Submit = function() {
		if ($scope.inProgress) {
			return;
		}
		$scope.$parent.page.errors = [];
		if (!$scope.IsPhoneValid()) {
			$scope.$parent.page.errors.push("Phone number is invalid");
		}
		if (!$scope.IsBarcodeValid()) {
			$scope.$parent.page.errors.push('Barcode is invalid');
		}
		if ($scope.$parent.page.errors.length) {
			return;
		}
		$scope.inProgress = true;
		$http({
			method: "POST",
			url: "/api/book/borrow",
			params: $scope.borrow
		}).success(function(data, status){
			// TODO: This seems not working. Fix it later
			$scope.$parent.page.messages = ["Borrow book successful"];
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
	$scope.page = 1;
	$scope.pages = [];
	$scope.SetQuery = function(query) {
		$scope.query = query;
	};
	$scope.Search = function() {
		var req;
		if ($scope.query) {
			req = $http({
				method: "GET",
				url: "/api/book/search",
				params: {"q": $scope.query, "page": $scope.page}
			});
		} else {
			req = $http({
				method: "GET",
				url: "/api/book/list",
				params: {"page": $scope.page}
			});
		}
		req.success(function(data, status){
			$scope.books = data.Books;
			$scope.page = data.Page;
			$scope.pages = range(1, data.NumPage+1);
			$scope.error = "";
		}).error(function(data, status){
			$scope.books = [];
			$scope.page = 1;
			$scope.pages = [];
			$scope.error = data;
		});
	};
	$scope.SetPage = function(page) {
		$scope.page = page;
		$scope.Search();
	};
	$scope.NextPage = function() {
		if ($scope.page < $scope.pages.length) {
			$scope.page++;
			$scope.Search();
		}
	};
	$scope.PrevPage = function() {
		if ($scope.page > 1) {
			$scope.page--;
			$scope.Search();
		}
	};
}]);

hoc5App.controller('BorrowerSearchCtrl', ["$scope", "$http", function($scope, $http){
	$scope.borrowers = [];
	$scope.query = "";
	$scope.error = "";
	$scope.page = 1;
	$scope.pages = [];
	$scope.SetQuery = function(query) {
		$scope.query = query;
	};
	$scope.CanBorrow = function(borrower) {
		// TODO: Implement it
		return true;
	};
	$scope.Search = function() {
		var req;
		if ($scope.query) {
			req = $http({
				method: "GET",
				url: "/api/borrower/search",
				params: {"q": $scope.query, "page": $scope.page}
			});
		} else {
			req = $http({
				method: "GET",
				url: "/api/borrower/list",
				params: {"page": $scope.page}
			});
		}
		req.success(function(data, status){
			$scope.borrowers = data.Borrowers;
			$scope.page = data.Page;
			$scope.pages = range(1, data.NumPage+1);
			$scope.error = "";
		}).error(function(data, status){
			$scope.borrowers = [];
			$scope.page = 1;
			$scope.pages = [];
			$scope.error = data;
		});
	};
	$scope.SetPage = function(page) {
		$scope.page = page;
		$scope.Search();
	};
	$scope.NextPage = function() {
		if ($scope.page < $scope.pages.length) {
			$scope.page++;
			$scope.Search();
		}
	};
	$scope.PrevPage = function() {
		if ($scope.page > 1) {
			$scope.page--;
			$scope.Search();
		}
	};
}]);