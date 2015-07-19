hoc5App.controller('BorrowerEditCtrl', [
	'$scope', '$routeParams', '$http', '$location', function($scope, $routeParams, $http, $location){
	$scope.$parent.page = {
		title: "Edit Borrower"
	};
	$scope.phone = $routeParams.phone;
	$scope.borrower = null;
	$http({
		method: "GET",
		url: "/api/borrower/" + $scope.phone
	}).success(function(data, status){
		$scope.borrower = data.Borrower;
		$scope.orignal = angular.copy($scope.borrower);
	}).error(function(data, status){
		$scope.$parent.page.errors = [data];
	});
	$scope.ValidateName = function(chineseName, englishName) {
		chineseName = chineseName || "";
		englishName = englishName || "";
		if (chineseName.length === 0 && englishName.length === 0) {
			return false;
		}
		return true;
	};
	$scope.Reset = function() {
		$scope.borrower = angular.copy($scope.orignal);
	};
	$scope.inProgress = false;
	$scope.Submit = function() {
		if ($scope.inProgress) {
			return;
		}
		$scope.inProgress = true;
		$http({
			method: "POST",
			url: "/api/borrower/" + $scope.phone,
			data: $scope.borrower
		}).success(function(data, status){
			$location.path("/menu");
		}).error(function(data, status){
			$scope.$parent.page.errors = [data];
			$scope.inProgress = false;
		});
	};
}]);

hoc5App.controller('BorrowerNewCtrl', [
	'$scope', '$http', '$location', function($scope, $http, $location){
	$scope.$parent.page = {
		title: "New Borrower"
	};
	$scope.ValidateName = function(chineseName, englishName) {
		chineseName = chineseName || "";
		englishName = englishName || "";
		if (chineseName.length === 0 && englishName.length === 0) {
			return false;
		}
		return true;
	};
	$scope.borrower = {};
	$scope.Reset = function() {
		$scope.borrower = {};
	};
	$scope.inProgress = false;
	$scope.Submit = function() {
		if ($scope.inProgress) {
			return;
		}
		$scope.inProgress = true;
		$http({
			method: "POST",
			url: "/api/borrower",
			data: $scope.borrower
		}).success(function(data, status){
			$location.path("/menu");
		}).error(function(data, status){
			$scope.$parent.page.errors = [data];
			$scope.inProgress = false;
		});
	};
}]);