function postToURL(path, params, method) {
    method = method || "post"; // Set method to post by default if not specified.

    // The rest of this code assumes you are not using a library.
    // It can be made less wordy if you use one.
    var form = document.createElement("form");
    form.setAttribute("method", method);
    form.setAttribute("action", path);
    form.setAttribute("target", "_blank");
    form.setAttribute("style", "display: none");

    for(var key in params) {
        if(params.hasOwnProperty(key)) {
            var hiddenField = document.createElement("input");
            hiddenField.setAttribute("type", "hidden");
            hiddenField.setAttribute("name", key);
            hiddenField.setAttribute("value", params[key]);

            form.appendChild(hiddenField);
         }
    }

    document.body.appendChild(form);
    form.submit();
    //document.body.removeChild(form);
}


hoc5App.controller('GenBarcodeCtrl', ['$scope', '$http', function($scope, $http){
	$scope.$parent.page = {};

	var initInput = {
		start: "000000000000",
		count: 0,
		step: 10
	};
	$scope.input = initInput;

	var result = [];
	var LayoutBarcode = function(data, countPerRow) {
		result.length = 0;
		for (var i = 0; i < data.length; i+=countPerRow) {
			var row = [];
			for (var j = i; j < data.length && j < i+countPerRow; j++) {
				row.push(data[j]);
			}
			result.push(row);
		}
		return result;	
	};
	$scope.Generate = function() {
		$http({
			"method": "GET",
			"url": "/barcode/book/" + $scope.input.start + "/" + $scope.input.count + "/" + $scope.input.step
		}).success(function(data, status){
			$scope.barcodes = $scope.barcodes || [];
			$.merge($scope.barcodes, data);
			$scope.layout = LayoutBarcode($scope.barcodes, 5);
			$scope.input = initInput;
		}).error(function(data, status){
		});
	};
	$scope.Remove = function(index) {
		$scope.barcodes.splice(index, 1);
		$scope.layout = LayoutBarcode($scope.barcodes, 5);
	}
	$scope.Print = function() {
		postToURL("/barcode/book/print.pdf", {
			barcodes: $scope.barcodes.join(",")
		})
	};
}]);
