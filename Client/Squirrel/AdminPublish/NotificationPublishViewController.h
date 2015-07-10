//
//  NotificationPublishViewController.h
//  Squirrel
//
//  Created by JamesMao on 4/12/15.
//  Copyright (c) 2015 JamesMao. All rights reserved.
//

#import <UIKit/UIKit.h>
#import "DataModel/AdminPublishViewData.h"
#import "AFHttpViewController/AFHttpViewController.h"
typedef enum
{
    EditType_View,
    EditType_Editable,
    EditType_Add,
    EditType_Delete,
} EditType;


@class UserRegisterDataItem;

@interface NotificationPublishViewController : AFHttpViewController

- (void)setEditType:(EditType)type;

- (void)setNotificationDataItem:(NotificationDataItem*)item;

@end
