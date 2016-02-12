var angular = require('angular');

require('./services.js');
require('./truncate.js');
require('./templates.js');

(function() {
  angular.module('ghca', [
    'ghca.services',
    'truncate',
    'templates',
    require('angular-moment'),
    require('angular-ui-bootstrap'),
    require('angular-ui-router')
  ]);
})();

require('./routes.js');
require('./error-handler.js');
require('./bsod.controller.js');
require('./search.controller.js');
require('./user.controller.js');
require('./stats.controller.js');
