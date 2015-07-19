hoc5App.controller('BookEditCtrl', [
	'$scope', '$routeParams', '$http', '$location', function($scope, $routeParams, $http, $location){
	$scope.$parent.page = {
		title: "Edit Book"
	};
	$scope.barcode = $routeParams.barcode;
	$scope.book = null;
	$http({
		method: "GET",
		url: "/api/book/" + $scope.barcode
	}).success(function(data, status){
		$scope.book = data.Book;
		$scope.orignal = angular.copy($scope.book);
		var d = new Date($scope.book.PublishDate || "");
		$scope.pubMonth = d.getUTCMonth()+1;
		$scope.pubYear = d.getUTCFullYear();
	}).error(function(data, status){
		$scope.$parent.page.errors = [data];
	});
	$scope.Reset = function() {
		$scope.book = angular.copy($scope.orignal);
		var d = new Date($scope.book.PublishDate || "");
		$scope.pubMonth = d.getUTCMonth()+1;
		$scope.pubYear = d.getUTCFullYear();
	};
	$scope.inProgress = false;
	$scope.Submit = function() {
		if ($scope.inProgress) {
			return;
		}
		if ($scope.editForm.$invalid) {
			return;
		}
		$scope.inProgress = true;
		$http({
			method: "POST",
			url: "/api/book/" + $scope.barcode,
			data: $scope.book
		}).success(function(data, status){
			$location.path("/menu");
		}).error(function(data, status){
			$scope.$parent.page.errors = [data];
			$scope.inProgress = false;
		});
	};
	$scope.ValidatePubDate = function(month, year) {
		if (!month && !year) {
			return true;
		}
		var m = parseInt(month);
		var y = parseInt(year);
		return (0 < m && m <= 12) && (y >= 1000 && y <= 2100);
	};
	$scope.UpdatePubDate = function() {
		if (!$scope.ValidatePubDate($scope.pubMonth, $scope.pubYear)) {
			return;
		}
		if (!$scope.pubMonth && !$scope.pubYear) {
			delete $scope.book.PublishDate;
		} else {
			var d = new Date(0);
			d.setUTCFullYear($scope.pubYear);
			d.setUTCMonth($scope.pubMonth-1);
			$scope.book.PublishDate = d;
		}
	};
}]);

hoc5App.controller('BookNewCtrl', [
	'$scope', '$http', '$location', function($scope, $http, $location){
	$scope.$parent.page = {
		title: "New Book"
	};
	$scope.book = {};
	$scope.Reset = function() {
		$scope.book = {};
		$scope.pubMonth = "";
		$scope.pubYear = "";
	};
	$scope.inProgress = false;
	$scope.Submit = function() {
		if ($scope.inProgress) {
			return;
		}
		if ($scope.editForm.$invalid) {
			return;
		}
		$scope.inProgress = true;
		$http({
			method: "POST",
			url: "/api/book",
			data: $scope.book
		}).success(function(data, status){
			$location.path("/menu");
		}).error(function(data, status){
			$scope.$parent.page.errors = [data];
			$scope.inProgress = false;
		});
	};
	$scope.ValidatePubDate = function(month, year) {
		if (!month && !year) {
			return true;
		}
		var m = parseInt(month);
		var y = parseInt(year);
		return (0 < m && m <= 12) && (y >= 1000 && y <= 2100);
	};
	$scope.UpdatePubDate = function() {
		if (!$scope.ValidatePubDate($scope.pubMonth, $scope.pubYear)) {
			return;
		}
		if (!$scope.pubMonth && !$scope.pubYear) {
			delete $scope.book.PublishDate;
		} else {
			var d = new Date(0);
			d.setUTCFullYear($scope.pubYear);
			d.setUTCMonth($scope.pubMonth-1);
			$scope.book.PublishDate = d;
		}
	};
}]);
