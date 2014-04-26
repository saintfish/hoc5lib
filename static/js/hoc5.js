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

hoc5App.controller('BookBorrowCtrl', ['$scope', function($scope){
	$scope.$parent.page = {
		title: "Borrow Books"
	};
	$scope.returnDate = function(){
		var d = new Date();
		d.setDate(d.getDate() + config.borrowPeriodDays);
		return d;
	}();
}]);

hoc5App.controller('SearchCtrl', ["$scope", "$http", function($scope, $http){
	$scope.books = [];
	$scope.query = "";
	$scope.error = "";
	$scope.page = 1;
	$scope.pages = [];
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