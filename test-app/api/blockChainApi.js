var express = require('express');

var bodyParser = require('body-parser');
var router = express.Router();
var app = express();
var chainControl = require('../utils/controller');

// var jwt = require('jsonwebtoken');
// var bcrypt = require('bcryptjs');

// for parsing application/json
router.use(bodyParser.json());

// for parsing application/xwww-
router.use(bodyParser.urlencoded({ extended: true }));

router.get('/get_all_data', function (req, res) {
    // User.find({}, { password: 0 }, function (err, users) {
    //     res.json(users);
    // });
    console.log('get_all_data called ');
    chainControl.read_all_data(req,res);
});

router.get('/domain_list', function (req, res) {
    // User.find({}, { password: 0 }, function (err, users) {
    //     res.json(users);
    // });
    console.log('get domain list called');
    chainControl.get_all_domain(req,res);
});

router.get('/company_list', function (req, res) {
    // User.find({}, { password: 0 }, function (err, users) {
    //     res.json(users);
    // });
    console.log('get company list called');
    chainControl.get_all_company(req,res);
});

router.get('/employee_list', function (req, res) {
    // User.find({}, { password: 0 }, function (err, users) {
    //     res.json(users);
    // });
    console.log('employee_list called');
    chainControl.get_all_employee(req,res);
});

router.get('/employee_details_by_id/:id', function (req, res) {
    // User.find({}, { password: 0 }, function (err, users) {
    //     res.json(users);
    // });
    console.log('employee_details_by_id called');
    chainControl.get_employee_by_id(req,res);
});

router.get('/delete_employee_by_id/:id/:domainName', function (req, res) {
    // User.find({}, { password: 0 }, function (err, users) {
    //     res.json(users);
    // });
    console.log('employee_details_by_id called');
    chainControl.delete_employee(req,res);
});

router.get('/disable_domain_by_id/:id/:domainName', function (req, res) {
    // User.find({}, { password: 0 }, function (err, users) {
    //     res.json(users);
    // });
    console.log('disable_domain_by_id called');
    chainControl.disable_domain(req,res);
});

router.get('/disable_company_by_id/:id/:companyName', function (req, res) {
    // User.find({}, { password: 0 }, function (err, users) {
    //     res.json(users);
    // });
    console.log('disable_company_by_id called');
    chainControl.disable_company(req,res);
});


router.post('/add_domain', function (req, res) {

    console.log('add domain called');
    chainControl.add_domain(req,res);

});

router.post('/add_company', function (req, res) {

    console.log('add company called');
    chainControl.add_company(req,res);


});

router.post('/add_employee', function (req, res) {
    console.log('add employee called');
    chainControl.add_employee(req,res);
});

router.post('/change_employee_domain', function (req, res) {

    console.log(' set_domain called');
    chainControl.set_domain(req,res);
    // req.body.password

    // User.findOne({ name: req.body.name }, function (err, user) {
    //     if (err) return res.status(500).send('Error on the server.');
    //     if (!user) return res.status(404).send('No user found.');

    //     var passwordIsValid = bcrypt.compareSync(req.body.password, user.password);
    //     if (!passwordIsValid) return res.status(401).send({ auth: false, token: null });

    //     var token = jwt.sign({ id: user._id }, config.secret, {
    //         expiresIn: '365d' // expires in 24 hours
    //     });

    //     res.status(200).send({ auth: true, token: token });
    // });

});

router.post('/change_employee_company', function (req, res) {

    console.log('set_company called');
    chainControl.set_company(req,res);
    // req.body.password

    // User.findOne({ name: req.body.name }, function (err, user) {
    //     if (err) return res.status(500).send('Error on the server.');
    //     if (!user) return res.status(404).send('No user found.');

    //     var passwordIsValid = bcrypt.compareSync(req.body.password, user.password);
    //     if (!passwordIsValid) return res.status(401).send({ auth: false, token: null });

    //     var token = jwt.sign({ id: user._id }, config.secret, {
    //         expiresIn: '365d' // expires in 24 hours
    //     });

    //     res.status(200).send({ auth: true, token: token });
    // });

});

module.exports = router;