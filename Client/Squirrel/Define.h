//
//  Define.h
//  Squirrel
//
//  Created by JamesMao on 6/7/15.
//  Copyright (c) 2015 JamesMao. All rights reserved.
//

#import <Foundation/Foundation.h>

#define BackGroundQueue dispatch_get_global_queue(DISPATCH_QUEUE_PRIORITY_DEFAULT, 0)


#define SERVER_IP_KEY @"SERVER_IP_KEY"
#define SERVER_IP_DEFAULT @"http://127.0.0.1:10443"


#define SERVER_IP [[NSUserDefaults standardUserDefaults] objectForKey:SERVER_IP_KEY]


#define URL_PART_DELETE_ONE_CLASS @"/API1/Classes/Delete/"
#define URL_PART_REGISTER_ONE_CLASS @"/API1/Classes/Register"
#define URL_PART_UNREGISTER_ONE_CLASS @"/API1/Classes/UnRegister"
#define URL_PART_QUERY_REGISTER_ONE_CLASS @"/API1/Classes/QueryRegisterStatus"
#define URL_PART_ONE_USER_CLASS @"/API1/Classes/User/"
#define URL_PART_CLASS_LIST @"/API1/Classes"
#define URL_PART_CLASS_POST @"/API1/Post"
#define URL_PART_ONE_CLASS @"/API1/OneClass/"
#define URL_PART_ONE_CLASS_USER_COUNT @"/API1/OneClassUserCount/"
#define URL_PART_ONE_CLASS_USERS @"/API1/OneClassUsers/"


#define LIGIN_VIEW_TAG 1000


#define NAVIGATION_BAR_BACKGROUND @"TabBarBackground"



