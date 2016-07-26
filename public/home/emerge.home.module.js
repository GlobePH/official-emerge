// home module
angular.module('emergeHome', [
    'uiGmapgoogle-maps'])
.config(['uiGmapGoogleMapApiProvider',
    function(uiGmapGoogleMapApiProvider) {
      /** Google Maps initialization **/
      uiGmapGoogleMapApiProvider.configure({
        // v: '3.20',
        libraries: 'weather, geometry, visualization'
      });
}]);
