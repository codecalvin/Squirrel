//
//  UserRegisterDataItem.h
//  Squirrel
//
//  Created by JamesMao on 6/16/15.
//  Copyright (c) 2015 JamesMao. All rights reserved.
//

#import <Foundation/Foundation.h>
#import "NotificationDataItem.h"

@interface UserRegisterDataItem : NSObject

- (void)setDataItemEelements:(NSMutableDictionary*)data;
- (NSMutableDictionary*)getDataItemEelements;

- (void)setUniqueKey:(NSString*)uniqueKey;
- (NSString*)getUniqueKey;
+ (NSString*)getUniqueKey:(NSMutableDictionary*)data;

- (void)setClassName:(NSString*)name;
- (NSString*)getClassName;
+ (NSString*)getClassName:(NSMutableDictionary*)data;

- (void)setUserUniqueKey:(NSString*)uniqueKey;
- (NSString*)getUserUniqueKey;
+ (NSString*)getUserUniqueKey:(NSMutableDictionary*)data;

- (void)setUserName:(NSString*)name;
- (NSString*)getUserName;
+ (NSString*)getUserName:(NSMutableDictionary*)data;


@end
