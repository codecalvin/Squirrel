//
//  Define.h
//  Squirrel
//
//  Created by JamesMao on 6/7/15.
//  Copyright (c) 2015 JamesMao. All rights reserved.
//

#import <Foundation/Foundation.h>

#define BackGroundQueue dispatch_get_global_queue(DISPATCH_QUEUE_PRIORITY_DEFAULT, 0) //1

#define SERVER_IP_KEY @"SERVER_IP_KEY"
#define SERVER_IP_DEFAULT @"http://192.168.1.103:5000"

#define SERVER_IP [[NSUserDefaults standardUserDefaults] objectForKey:SERVER_IP_KEY]

#define URL_PART_DELETE_ONE_CLASS @"/API1/Classes/Delete/"
#define URL_PART_REGISTER_ONE_CLASS @"/API1/Classes/Register"
#define URL_PART_ONE_USER_CLASS @"/API1/Classes/User/"

#define ClassListURL [NSURL URLWithString: [NSString stringWithFormat:@"%@/API1/Classes", SERVER_IP]]

#define OneClassURLBase [NSString stringWithFormat:@"%@/API1/OneClass/", SERVER_IP]

