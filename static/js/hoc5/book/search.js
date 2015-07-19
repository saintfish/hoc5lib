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