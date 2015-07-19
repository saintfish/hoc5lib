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
