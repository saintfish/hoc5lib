hoc5App.controller('BorrowerManageCtrl', [
	'$scope', '$location', function($scope, $location){
	$scope.$parent.page = {
		title: "Manage Borrowers"
	};
	$scope.Edit = function(borrower) {
		$location.path("/borrower/" + borrower.Phone + "/edit");
	};
	$scope.New = function() {
		$location.path("/borrower/new");
	};
}]);
