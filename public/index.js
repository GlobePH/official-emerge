'use strict';

var emergeApp = angular.module('emergeApp', [
    'ngRoute',
    'mobile-angular-ui',
    'uiGmapgoogle-maps'
]);

// TODO: Search for $transform
// emergeApp.run(function($transform) {
//   window.$transform = $tranform;
// });

emergeApp.config(['$routeProvider', 'uiGmapGoogleMapApiProvider',
    function($routeProvider, uiGmapGoogleMapApiProvider) {

  $routeProvider.when('/', {
    templateUrl:      'home.html',
    reloadOnSearch:   false
  });

  /** Google Maps initialization **/
  uiGmapGoogleMapApiProvider.configure({
    // v: '3.20',
    // libraries: 'weather, geometry, visualization'
  });

}]);

emergeApp.controller('MainController',
    [ '$rootScope',
      '$scope',
      'uiGmapGoogleMapApi',
    function($rootScope, $scope, uiGmapGoogleMapApi) {

      uiGmapGoogleMapApi.then(function(maps) {
        $scope.googleVersion = maps.version;
        maps.visualRefresh = true;
      });

      $scope.map = {
        center: {
          latitude: 45,
          longitude: -73
        },
        zoom: 8
      };

      $scope.sample = 'sample';

}]);

/*
 * Database: Firebase
 * Url: https://sag1application-fe490.firebaseio.com/members 
 * Data strucure:
 * Training Document object:
 * {
 *  date: dateToday,
 *  attendance: [
 *    {fName: firstName, lName: lastName, paid: true},
 *    {fName: firstName, lName: lastName, paid: false},
 *  ],
 *  totalFeesCollected: 12345
 * }
 *
 * Member Document object:
 * {
 *  fName           : firstName.
 *  lName           : lastName,
 *  position        : bowOrStroke,
 *  memberSince     : dateOfMembership,
 *  totalDues       : runningTotalDue
 * }
 *
 * Non member Document object:
 * {
 *  fName           : firstName.
 *  lName           : lastName,
 *  position        : bowOrStroke,
 * }
 *
 */
