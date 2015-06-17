//
//  UserViewData.h
//  Squirrel
//
//  Created by JamesMao on 6/7/15.
//  Copyright (c) 2015 JamesMao. All rights reserved.
//

#import <Foundation/Foundation.h>
#import "AdminPublishViewData.h"

@interface UserViewData : NSObject

+ (UserViewData*)singleton;

- (int)getNotificationCount;
- (NSMutableDictionary*)getNotificationItem:(int)index;

- (BOOL)setNotificationDataItem:(NotificationDataItem*)dataItem;
- (void)reset;
- (void)save;

- (void)setJsonDictionaryData:(NSDictionary*)jsonDictionary;

@end

