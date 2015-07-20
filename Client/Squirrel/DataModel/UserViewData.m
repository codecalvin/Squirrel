//
//  UserViewData.m
//  Squirrel
//
//  Created by JamesMao on 6/7/15.
//  Copyright (c) 2015 JamesMao. All rights reserved.
//

#import "UserViewData.h"

#define UserNotificationDataFileName @"UserNotificationDataFileName.plist"

@interface UserViewData ()
{
    NSMutableArray* notificationDictionary_;
    NotificationDataItem* notificationDataItem_;
}

- (BOOL)deleteNotificationDataItem:(NotificationDataItem*)dataItem;

@end

@implementation UserViewData

+ (UserViewData*)singleton
{
    static UserViewData* theUserViewDatanstance;
    if (theUserViewDatanstance == nil)
    {
        theUserViewDatanstance =  [[UserViewData alloc] init];
    }
    return theUserViewDatanstance;
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
    notificationDataItem_ = [[NotificationDataItem alloc] init];
    
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
    NSString* dataFilePathFullName = [documentsDirectory stringByAppendingPathComponent:UserNotificationDataFileName];
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

- (void)reset
{
    [notificationDictionary_ removeAllObjects];
}

- (void)save
{
    NSString* filePath = [self dataFilePath];
    if ([notificationDictionary_ writeToFile:filePath atomically:YES] == NO)
    {
        [notificationDictionary_ writeToFile:filePath atomically:NO];
    }
}

- (void)setJsonDictionaryData:(NSDictionary*)jsonDictionary
{
    NSArray* keys = [jsonDictionary allKeys];
    for (NSString* key in keys)
    {
        NSMutableDictionary* valueDictionaryItem = [jsonDictionary objectForKey:key];
        [notificationDataItem_ setDataItemEelements:valueDictionaryItem];
        [self setNotificationDataItem:notificationDataItem_];
    }
}


@end
