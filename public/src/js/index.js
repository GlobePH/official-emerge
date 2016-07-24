'use strict';

var emergeApp = angular.module('emergeApp', [
    'ngRoute',
    'mobile-angular-ui',
    'uiGmapgoogle-maps',
    'ngWebSocket'
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

// Socket.IO factory/service
// emergeApp.factory('mySocket', function(socketFactory) {
//   return socketFactory();
//   var domainUrl = 'https://localhost:3001';
//   // var domainUrl = 'https://emerge-app.herokuapp.com';
//   var listen = 'https://echo.websocket.org';
//   var listenUrl = '/api/channel'
//   var myIoSocket = io.connect(listen,
//     {
//       secure: true,
// 
//     });
// 
//   var mySocket = socketFactory({
//     ioSocket: myIoSocket
//   });
// 
//   return mySocket;
// });

emergeApp.factory('socket', function($websocket) {
  var url = 'ws://echo.websocket.org';
  // var url = 'https://emerge-app.herokuapp.com/api/channel';
  var dataStream = $websocket(url);
  var collection = [];

  dataStream.onMessage(function(message) {
    collection.push(JSON.parse(message.data));
    console.log( 'data received is ' + JSON.stringify(JSON.parse(message.data)));
  });

  var methods = {
    collection: collection,
    get: function() {
      dataStream.send(JSON.stringify({action: 'get'}));
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

      /** Markers **/
      // $scope.markers = [];
      // $scope.markerCount = 0;
      // $scope.markers.push(plot(14, 122));
      
      $scope.map = {
        center: {
          latitude: 13,
          longitude: 122
        },
        zoom: 7
      };


      /** Socket listeners to other servers **/
      // mySocket.on('connect', function() {
      //   console.log('connected...');
      // });

      // mySocket.on('disconnect', function() {
      //   console.log('disconnected...');
      // });

      // mySocket.on('message', function(data) {
      //   console.log('message received');
      //   console.log('data is ' + data.hello);

      //   mySocket.emit('front', { hello: 'front' });
      //   console.log('front is sent to backend...');
      // });

      // socket.send('Data from emerge sent to echo.websocket! Must show in console!');
      socket.get();

}]);

