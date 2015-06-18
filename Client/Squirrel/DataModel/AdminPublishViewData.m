//
//  AdminPublishViewData.m
//  Squirrel
//
//  Created by JamesMao on 4/11/15.
//  Copyright (c) 2015 JamesMao. All rights reserved.
//

#import "AdminPublishViewData.h"

#define PublishNotificationDataFileName @"PublishNotificationDataFileName.plist"

@interface AdminPublishViewData()
{
    NSMutableArray* notificationDictionary_;
}

- (int)getNotificationCount;
- (NSMutableDictionary*)getNotificationItem:(int)index;

@end

@implementation AdminPublishViewData

+ (AdminPublishViewData*)singleton
{
    static AdminPublishViewData* theAdminPublishViewDataInstance;
    if (theAdminPublishViewDataInstance == nil)
    {
        theAdminPublishViewDataInstance =  [[AdminPublishViewData alloc] init];
    }
    return theAdminPublishViewDataInstance;
}

- (id)init
{
    if (self = [super init])
    {
        [self initDataDictionaries];
    }
    return self;
}

- (void)initDataDictionaries
{
    NSString* filePath = [self dataFilePath];
    if ([[NSFileManager defaultManager] fileExistsAtPath:filePath])
    {
        notificationDictionary_ = [[NSMutableArray alloc]initWithContentsOfFile:filePath];
    }
    else
    {
        notificationDictionary_ = [[NSMutableArray alloc] init];
    }
}

- (NSString* )dataFilePath
{
    NSArray* path = NSSearchPathForDirectoriesInDomains(NSDocumentDirectory, NSUserDomainMask, YES);
    NSString* documentsDirectory = [path objectAtIndex:0];
    NSString* dataFilePathFullName = [documentsDirectory stringByAppendingPathComponent:PublishNotificationDataFileName];
    return dataFilePathFullName;
}

- (int)getNotificationCount
{
    return (int)[notificationDictionary_ count];
}

- (NSMutableDictionary*)getNotificationItem:(int)index
{
    if (0 <= index && index < (int)[notificationDictionary_ count])
    {
        return [notificationDictionary_ objectAtIndex:index];
    }
    
    return nil;
}

- (BOOL)setNotificationDataItem:(NotificationDataItem*)dataItem
{
    for (int index = 0; index < [notificationDictionary_ count]; index++)
    {
        NSMutableDictionary*  dataItemEelements = [self getNotificationItem:index];
        if (dataItemEelements == nil)
        {
            continue;
        }
        
        NSString* uniqueKeyCandidate = [NotificationDataItem getUniqueKey:dataItemEelements];
        NSString* uniqueKey = [dataItem getUniqueKey];
        if ([uniqueKey compare:uniqueKeyCandidate] == NSOrderedSame)
        {
            [notificationDictionary_ removeObject:dataItemEelements];
            [notificationDictionary_ addObject:[NSMutableDictionary dictionaryWithDictionary:[dataItem getDataItemEelements] ]];
            return YES;
        }
    }
    
    NSMutableDictionary* dataItemEelements = [[NSMutableDictionary alloc] initWithDictionary:[dataItem getDataItemEelements]];
    
    [notificationDictionary_ addObject:dataItemEelements];
    
    [self save];
    return YES;
}

- (BOOL)deleteNotificationDataItem:(NotificationDataItem*)dataItem;
{
    for (int index = 0; index < [notificationDictionary_ count]; index++)
    {
        NSMutableDictionary*  dataItemEelements = [self getNotificationItem:index];
        if (dataItemEelements == nil)
        {
            continue;
        }
        
        NSString* uniqueKeyCandidate = [NotificationDataItem getUniqueKey:dataItemEelements];
        NSString* uniqueKey = [dataItem getUniqueKey];
        if ([uniqueKey compare:uniqueKeyCandidate] == NSOrderedSame)
        {
            [notificationDictionary_ removeObject:dataItemEelements];
            return YES;
        }
    }
    
    return NO;
}

- (void)save
{
    NSString* filePath = [self dataFilePath];
    if ([notificationDictionary_ writeToFile:filePath atomically:YES] == NO)
    {
        [notificationDictionary_ writeToFile:filePath atomically:NO];
    }
}

@end
