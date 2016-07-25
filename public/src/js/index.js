'use strict';

var emergeApp = angular.module('emergeApp', [
    'ngRoute',
    'mobile-angular-ui',
    'uiGmapgoogle-maps',
    'ngWebSocket',
    'emergeSidebar'
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
    libraries: 'weather, geometry, visualization'
  });

}]);

// Websocket provider
emergeApp.factory('socket', function($websocket, $window) {
  // var url = 'ws://echo.websocket.org';
  var url = 'wss://emerge-app.herokuapp.com/api/channel';
  var dataStream = $websocket(url);
  var collection = [];

  dataStream.onMessage(function(message) {
    collection.push(JSON.parse(message.data));
    console.log( 'data received is ' +
      JSON.stringify(JSON.parse(message.data)));
    $window.alert('data received from websocket server: ' + 
      JSON.stringify(JSON.parse(message.data)));
  });

  var methods = {
    collection: collection,
    get: function() {
      dataStream.send(JSON.stringify({action: 'get'}));
    },
    // Send a sync request to the server
    sync: function() {
      ;
    }
  };
  return methods;
});


emergeApp.controller('MainController',
    [ '$rootScope',
      '$scope',
      'uiGmapGoogleMapApi',
      'socket',
      // 'plot',
    function($rootScope, $scope, uiGmapGoogleMapApi, socket) {

      uiGmapGoogleMapApi.then(function(maps) {
        // $scope.googleVersion = maps.version;
        // maps.visualRefresh = true;
      });

      /** GMaps Markers TODO! make markers work **/
      // $scope.markers = [];
      // $scope.markerCount = 0;
      // $scope.markers.push(plot(14, 122));
      
      // Google maps demonstration
      $scope.map = {
        center: {
          latitude: 13,
          longitude: 122
        },
        zoom: 7
      };

      // Web socket demonstration
      socket.get();

}]);

