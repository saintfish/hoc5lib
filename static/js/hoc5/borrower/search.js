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
