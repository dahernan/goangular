'use strict';

angular.module('clientApp')
  .controller('MainCtrl', function ($scope, $http) {
    $scope.awesomeThings = [
      'HTML5 Boilerplate',
      'AngularJS',
      'Karma'
    ];

    $scope.add = function(){
    	$http.get('/api/hello/' + $scope.something ).success(function(data){
    		$scope.awesomeThings.push(data);
    	});
    	
    }
  });
