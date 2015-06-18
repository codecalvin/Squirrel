//
//  AdminPublishViewData.h
//  Squirrel
//
//  Created by JamesMao on 4/11/15.
//  Copyright (c) 2015 JamesMao. All rights reserved.
//

#import <Foundation/Foundation.h>
#import "NotificationDataItem.h"


@interface AdminPublishViewData : NSObject

+ (AdminPublishViewData*)singleton;

- (BOOL)setNotificationDataItem:(NotificationDataItem*)dataItem;

@end
