function SudokuCtrl($scope, $http) {

	$scope.inputString = "000000068900000002000400500041000000000035000050000000000800010300000700000100400";
	$scope.unsolved = emptyBoard()
	$scope.solved = emptyBoard()
	$scope.solving = false

	$scope.solveBoard = function() {
		$scope.errorMessage = null;
		$scope.solved = emptyBoard();
		$scope.solving = true;

		var str = boardArrayToString($scope.unsolvedBoard());

		$http.get('/sudoku/' + str
		).success(function(data) {
			$scope.solved = data;
		}).error(function(data, status, headers, config) {
			if (status == 0) {
				$scope.errorMessage = "Failed to reach Sudoku solver"
			} else {
				$scope.errorMessage = "Error: " + data;
			}
			$scope.solved = null
		}).finally(function() {
			$scope.solving = false;
		});
	};

	$scope.unsolvedBoard = function() {
		return boardStringToArray($scope.inputString, $scope.unsolved);
	}
}

function emptyBoard() {
	var arr = [];
	for (var row=0; row<9; row++) {
		arr.push(new Array(9));
	}
	return arr;
}

function boardStringToArray(str, arr) {
	str = str.replace(/ /, "")
	for (var row=0; row<9; row++) {
		for (var col=0; col<9; col++) {
			var idx = row * 9 + col;
			var val = 0;
			if (idx < str.length) {
				var ch = str.charAt(idx)
				if (ch >= "0" && ch <= "9") {
					val = ch - "0";
				}
			}
			arr[row][col] = val;
		}
	}
	return arr;
}

function boardArrayToString(arr) {
	str = ""
	for (var row=0; row<9; row++) {
		for (var col=0; col<9; col++) {
			str += arr[row][col]
		}
	}
	return str;
}
