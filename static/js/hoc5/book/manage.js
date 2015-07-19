hoc5App.controller('BookManageCtrl', [
	'$scope', '$location', function($scope, $location){
	$scope.$parent.page = {
		title: "Manage Books"
	};
	$scope.Edit = function(book) {
		$location.path("/book/" + book.Barcode + "/edit");
	};
	$scope.New = function() {
		$location.path("/book/new");
	};
}]);
