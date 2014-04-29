var config = {
	"borrowPeriodDays": 15,
	"numBookOneCanBorrow": 3
};

var hoc5App = angular.module("hoc5App", ["ngRoute", "ui.bootstrap", "ui.utils"]);

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
		when("/book/return", {
			templateUrl: "partials/book/return.html",
			controller: "BookReturnCtrl"
		}).
		when("/book/manage", {
			templateUrl: "partials/book/manage.html",
			controller: "BookManageCtrl"
		}).
		when("/book/new", {
			templateUrl: "partials/book/edit.html",
			controller: "BookNewCtrl"
		}).
		when("/book/:barcode/edit", {
			templateUrl: "partials/book/edit.html",
			controller: "BookEditCtrl"
		}).
		when("/borrower/manage", {
			templateUrl: "partials/borrower/manage.html",
			controller: "BorrowerManageCtrl"
		}).
		when("/borrower/new", {
			templateUrl: "partials/borrower/edit.html",
			controller: "BorrowerNewCtrl"
		}).
		when("/borrower/:phone/edit", {
			templateUrl: "partials/borrower/edit.html",
			controller: "BorrowerEditCtrl"
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
		return (0 < m && m <= 12) && (y <= 2100);
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
		return (0 < m && m <= 12) && (y <= 2100);
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