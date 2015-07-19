function ListStock() {
	var key, keyObj;
	var result = [];
	for (var i = 0; i < localStorage.length; i++) {
		key = localStorage.key(i);
		keyObj = JSON.parse(key);
		if (keyObj.begin && keyObj.end) {
			result.push(keyObj);
		}
	}
	return result;
}

function NewStock(stock) {
	var key = JSON.stringify(stock);
	localStorage.setItem(key, "");
}

function DeleteStock(stock) {
	var key = JSON.stringify(stock);
	localStorage.removeItem(key);
}

function GetStock(stock) {
	var key = JSON.stringify(stock);
	var value = localStorage.getItem(key);
	if (value == "") {
		return null;
	} else {
		return JSON.parse(value);
	}
}

function SetStock(stock, books) {
	var key = JSON.stringify(stock);
	var value = JSON.stringify(books);
	localStorage.setItem(key, value);
}

hoc5App.controller('BookStockCtrl', [
	'$scope', '$location', function($scope, $location){
	$scope.$parent.page = {
		title: "Stock Books"
	};

	$scope.stocks = ListStock();
	$scope.NewStock = function(stock) {
		NewStock({'begin': stock.begin, 'end': stock.end});
		$scope.GotoStock(stock);
	};
	$scope.DeleteStock = function(stock) {
		DeleteStock({'begin': stock.begin, 'end': stock.end});
		$scope.stocks = ListStock();
	};
	$scope.GotoStock = function(stock) {
		$location.path("/book/stock/" + stock.begin + "/" + stock.end);
	};
}]);

hoc5App.controller('BookStockEntryCtrl', [
	'$scope', '$routeParams', '$http', function($scope, $routeParams, $http){
	$scope.stock = {
		'begin': $routeParams.begin,
		'end': $routeParams.end
	};
	$scope.$parent.page = {
		title: "Stock Books " + $scope.stock.begin + " - " + $scope.stock.end
	};
	$scope.books = GetStock($scope.stock);
	if (!$scope.books) {
		$http({
			"method": "GET",
			"url": "/api/book/" + $scope.stock.begin + "/" + $scope.stock.end
		}).success(function(data, status){
			$scope.books = {
				"borrowed": [],
				"uncategorized": [],
				"stock": []
			};
			$(data).each(function(key, book) {
				if (!book.Availability) {
					$scope.books.borrowed.push(book);
				} else {
					$scope.books.uncategorized.push(book);
				}
			});
			SetStock($scope.stock, $scope.books);
		});
	}
	$scope.ValidateBarcode = function(barcode) {
		for (var i = 0; i < $scope.books.uncategorized.length; i++) {
			if ($scope.books.uncategorized[i].Barcode === barcode) {
				return true;
			}
		}
		return false;
	};
	$scope.AddStock = function(barcode) {
		for (var i = 0; i < $scope.books.uncategorized.length; i++) {
			var book = $scope.books.uncategorized[i];
			if (book.Barcode === barcode) {
				$scope.books.stock.push(book);
				$scope.books.uncategorized.splice(i, 1);
				SetStock($scope.stock, $scope.books);
				$scope.barcode = "";
				return;
			}
		}
	};
}]);
