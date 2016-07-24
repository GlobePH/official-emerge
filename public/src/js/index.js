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

// Socket.IO factory/service
emergeApp.factory('mySocket', function(socketFactory) {
  // var listenUrl = 'https://emerge-app.herokuapp.com/#/api/channel'
  var domainUrl = 'https://emerge-app.herokuapp.com';
  var listenUrl = '/api/channel'
  var myIoSocket = io.connect(domainUrl + listenUrl, {secure: true});
  // io.configure(function() {
  //   io.set('transports', ['xhr-polling']);
  //   io.set('polling duration', 10);
  // });
  var mySocket = socketFactory({
    ioSocket: myIoSocket
  });

  return mySocket;
});

// Google plotter factory/service
// emergeApp.service('plot' ,function(markerCount, lat, lon) {
//   // $scope.markerCount++;
//   var marker = {
//     idKey: markerCount,
//     coord: {
//       latitude: lat,
//       longitude: lon
//     }
//   };
//   return marker;
// });


emergeApp.controller('MainController',
    [ '$rootScope',
      '$scope',
      'uiGmapGoogleMapApi',
      'mySocket',
      // 'plot',
    function($rootScope, $scope, uiGmapGoogleMapApi, mySocket) {

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
      mySocket.on('connection', function() {
        console.log('connected...');
      });

      mySocket.on('message', function(data) {
        console.log('message received');
        console.log('data is ' + data.hello);

        mySocket.emit('front', { hello: 'front' });
        console.log('front is sent to backend...');
      });

}]);

