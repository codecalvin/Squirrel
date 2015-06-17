//
//  MeViewData.h
//  Squirrel
//
//  Created by JamesMao on 6/10/15.
//  Copyright (c) 2015 JamesMao. All rights reserved.
//

#import <Foundation/Foundation.h>
#import "AdminPublishViewData.h"
#import "UserRegisterDataItem.h"


@interface MeViewData : NSObject

+ (MeViewData*)singleton;

@property (nonatomic, strong)NSString* userUniqueName;

- (int)getNotificationCount;
- (NSMutableDictionary*)getNotificationItem:(int)index;

- (BOOL)deleteNotificationDataItem:(NotificationDataItem*)dataItem;
- (void)reset;
- (void)save;

- (void)setJsonDictionaryData:(NSDictionary*)jsonDictionary;

@end
