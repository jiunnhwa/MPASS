package main

/*

/api/job/

*/

/*

Errors:

2021/07/18 06:19:47 Error 1461: Can't create more than max_prepared_stmt_count statements (current value: 16382)


Result = '{"code": 21606, "message": "The From phone number +6594497191 is not a valid, SMS-capable inbound phone number
or short code for your account.", "more_info": "https://www.twilio.com/docs/errors/21606", "status": 400}'

Result = '{"code": 21211, "message": "The \'To\' number +6598219019 is not a valid phone number.", "more_info": "https://www.twilio.com/docs/errors/21211", "status": 400}'


2021/08/15 05:28:12 AsJSON: 200 {"RID":147,"RecordsAffected":1} <nil>
error: 0        404 Not Found   Body : {"code": 20404, "message": "The requested resource /2010-04-01/Accounts//Messages.json was not found", "more_info": "https://www.twilio.com/docs/errors/20404", "status": 404}

Status: 401 Unauthorized        Body : {"code": 20003, "detail": "", "message": "Authenticate", "more_info": "https://www.twilio.com/docs/errors/20003", "status": 401}

*/
