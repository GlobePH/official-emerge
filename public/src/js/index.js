'use strict';

var emergeApp = angular.module('emergeApp', [
    'ngRoute',
    'mobile-angular-ui',
    'uiGmapgoogle-maps',
    'btford.socket-io'
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

emergeApp.factory('mySocket', function(socketFactory) {
  var myIoSocket = io.connect('http://localhost:3001');
  var mySocket = socketFactory({
    ioSocket: myIoSocket
  });

  return mySocket;
});

// Wrap Socket.io into an Angular service
// http://www.html5rocks.com/en/tutorials/frameworks/angular-websockets/
// emergeApp.factory('socket', function($rootScope) {
//   var listenUrl = ('http://localhost:3001');
//   var socket = io.connect(listenUrl);
//   return {
//     on: function(eventName, callback) {
//       socket.on(eventName, function() {
//         var args = arguments;
//         $rootScope.$apply(function() {
//           if (callback) {
//             callback.apply(socket, args);
//           }
//         });
//       });
//     },
//     emit: function(eventName, data, callback) {
//       socket.emit(eventName, data, function() {
//         var args = arguments;
//         $rootScope.$apply(function() {
//           if (callback) {
//             callback.apply(socket, args);
//           }
//         });
//       });
//     }
//   };
// });


emergeApp.controller('MainController',
    [ '$rootScope',
      '$scope',
      'uiGmapGoogleMapApi',
      'mySocket',
    function($rootScope, $scope, uiGmapGoogleMapApi, mySocket) {

      uiGmapGoogleMapApi.then(function(maps) {
        // $scope.googleVersion = maps.version;
        // maps.visualRefresh = true;
      });
      
      $scope.map = {
        center: {
          latitude: 13,
          longitude: 122
        },
        zoom: 1
      };

      $scope.sample = 'sample';

      /** Socket listeners to other servers **/
      mySocket.on('connection', function() {
        // $scope.alert('connected...');
        console.log('connected...');
      });

      mySocket.on('message', function(data) {
        // $scope.alert(data);
        console.log('message received');
        console.log('data is ' + data.hello);

        mySocket.emit('front', { hello: 'front' });
        console.log('front is sent to backend...');
      });

}]);

