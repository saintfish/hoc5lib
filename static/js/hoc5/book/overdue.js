hoc5App.controller('BookOverdueCtrl', [
	'$scope', '$http', function($scope, $http){
	$scope.$parent.page = {
		title: "Overdue Books"
	};
	$http({
		method: "GET",
		url: "/api/book/overdue"
	}).success(function(data, status){
		$scope.entry = data.Entry;
	}).error(function(data, status){
		$scope.$parent.page.errors = [data];
	});
	$scope.Print = function() {
		window.print();
	};
}]);