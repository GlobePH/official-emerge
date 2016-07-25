// emergeHome component
'use strict';
angular.module('emergeHome')
.component('emergeHomeMap', {
  templateUrl: 'home/emerge.home.template.html',
  controller: [
    '$rootScope',
    '$scope',
    'uiGmapGoogleMapApi',
  function homeController($rootScope, $scope, uiGmapGoogleMapApi) {
    /** Controller function for home **/

    uiGmapGoogleMapApi.then(function(maps) {
      // $scope.googleVersion = maps.version;
      // maps.visualRefresh = true;
    });

    $scope.map = {
      center: {
        latitude: 13,
         longitude: 122
      },
      zoom: 7
    };

  }]
});
