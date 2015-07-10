//
//  NotificationDataItem.h
//  Squirrel
//
//  Created by JamesMao on 6/16/15.
//  Copyright (c) 2015 JamesMao. All rights reserved.
//

#import <Foundation/Foundation.h>

typedef enum
{
    ElementType_ClassName,
    ElementType_UniqueKey,
    ElementType_ClassTime,
    ElementType_ClassTeacher,
    ElementType_ClassDescription,
    ElementType_ClassStudent,
    
    ElementType_UserName,
    ElementType_UserUniqueKey,
    ElementType_ClassRegisteredCount,
    
} ElementType;

//-------------------------------------------------------------------------------------
// class NotificationDataItem
//
@interface NotificationDataItem : NSObject

- (void)setDataItemEelements:(NSMutableDictionary*)data;
- (NSMutableDictionary*)getDataItemEelements;

- (void)setUniqueKey:(NSString*)uniqueKey;
- (NSString*)getUniqueKey;
+ (NSString*)getUniqueKey:(NSMutableDictionary*)data;

- (void)setClassName:(NSString*)name;
- (NSString*)getClassName;
+ (NSString*)getClassName:(NSMutableDictionary*)data;

- (void)setClassTime:(NSString*)name;
- (NSString*)getClassTime;
+ (NSString*)getClassTime:(NSMutableDictionary*)data;

- (void)setClassTeacher:(NSString*)name;
- (NSString*)getClassTeacher;
+ (NSString*)getClassTeacher:(NSMutableDictionary*)data;

- (void)setClassDescription:(NSString*)name;
- (NSString*)getClassDescription;
+ (NSString*)getClassDescription:(NSMutableDictionary*)data;

- (void)setClassMaxStudent:(NSString*)name;
- (NSString*)getClassMaxStudent;
+ (NSString*)getClassMaxStudent:(NSMutableDictionary*)data;

- (void)setNotificationDataItem:(NotificationDataItem*)item;

+ (NSString*)getCurrentTimeString;
+ (NSString*)getUniqueIndex;
+ (NSString*)getKey:(ElementType)type;

@end
