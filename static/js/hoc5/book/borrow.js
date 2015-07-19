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
			headers: {'Content-Type': 'application/x-www-form-urlencoded'},
			data: $.param($scope.borrow)
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
