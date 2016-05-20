var angular = require('angular');
var chart = require('angular-chart.js');

require('./services.js');
require('./truncate.js');
require('./templates.js');

(function() {
  angular.module('ghca', [
    'ghca.services',
    'truncate',
    'templates',
    'chart.js',
    require('angular-moment'),
    require('angular-ui-bootstrap'),
    require('angular-ui-router')
  ]).config(function($locationProvider) {
    $locationProvider.html5Mode(true).hashPrefix('!');
  })
  ;
})();


require('./routes.js');
require('./error-handler.js');
require('./bsod.controller.js');
require('./chart.controller.js');
require('./search.controller.js');
require('./user.controller.js');
require('./stats.controller.js');
