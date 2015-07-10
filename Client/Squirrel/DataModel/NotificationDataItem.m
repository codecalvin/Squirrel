//
//  NotificationDataItem.m
//  Squirrel
//
//  Created by JamesMao on 6/16/15.
//  Copyright (c) 2015 JamesMao. All rights reserved.
//

#import "NotificationDataItem.h"


//-------------------------------------------------------------------------------------
// class NotificationDataItem
//
@interface NotificationDataItem()
{
    NSMutableDictionary* dataItemEelements_;
}
+ (NSString*)getKey:(ElementType)type;

@end

@implementation NotificationDataItem

- (id)init
{
    if (self = [super init])
    {
        [self doInitialize];
    }
    return self;
}

- (void)doInitialize
{
    dataItemEelements_ = [[NSMutableDictionary alloc] init];
    
    NSString *dateString = [NotificationDataItem getCurrentTimeString];
    NSString * uniqueKey = [NSString stringWithFormat:@"%@_%@", dateString, [NotificationDataItem getUniqueIndex]];
    
    [self setUniqueKey:uniqueKey];
}

- (void)setDataItemEelements:(NSMutableDictionary*)data
{
    dataItemEelements_ = [[NSMutableDictionary alloc] initWithDictionary:data];
}

- (NSMutableDictionary*)getDataItemEelements
{
    return dataItemEelements_;
}

- (void)setUniqueKey:(NSString*)uniqueKey
{
    uniqueKey = [uniqueKey stringByReplacingOccurrencesOfString:@" " withString:@"_"];
    NSString* key = [NotificationDataItem getKey:ElementType_UniqueKey];
    [dataItemEelements_ setObject:uniqueKey forKey:key];
}

+ (NSString*)getUniqueKey:(NSMutableDictionary*)data
{
    NSString* key = [NotificationDataItem getKey:ElementType_UniqueKey];
    return [data objectForKey:key];
}

- (NSString*)getClassName
{
    NSString* key = [NotificationDataItem getKey:ElementType_ClassName];
    NSString* unique = [dataItemEelements_ objectForKey:key];
    return unique;
}

- (void)setClassName:(NSString*)value
{
    NSString* key = [NotificationDataItem getKey:ElementType_ClassName];
    [dataItemEelements_ setObject:value forKey:key];
}

+ (NSString*)getClassName:(NSMutableDictionary*)data
{
    NSString* key = [NotificationDataItem getKey:ElementType_ClassName];
    return [data objectForKey:key];
}

- (NSString*)getClassTime
{
    NSString* key = [NotificationDataItem getKey:ElementType_ClassTime];
    NSString* unique = [dataItemEelements_ objectForKey:key];
    return unique;
}

- (void)setClassTime:(NSString*)value
{
    NSString* key = [NotificationDataItem getKey:ElementType_ClassTime];
    [dataItemEelements_ setObject:value forKey:key];
}

+ (NSString*)getClassTime:(NSMutableDictionary*)data
{
    NSString* key = [NotificationDataItem getKey:ElementType_ClassTime];
    return [data objectForKey:key];
}

- (NSString*)getClassTeacher
{
    NSString* key = [NotificationDataItem getKey:ElementType_ClassTeacher];
    NSString* unique = [dataItemEelements_ objectForKey:key];
    return unique;
}

- (void)setClassTeacher:(NSString*)value
{
    NSString* key = [NotificationDataItem getKey:ElementType_ClassTeacher];
    [dataItemEelements_ setObject:value forKey:key];
}

+ (NSString*)getClassTeacher:(NSMutableDictionary*)data
{
    NSString* key = [NotificationDataItem getKey:ElementType_ClassTeacher];
    return [data objectForKey:key];
}

- (NSString*)getClassDescription
{
    NSString* key = [NotificationDataItem getKey:ElementType_ClassDescription];
    NSString* unique = [dataItemEelements_ objectForKey:key];
    return unique;
}

- (void)setClassDescription:(NSString*)value
{
    NSString* key = [NotificationDataItem getKey:ElementType_ClassDescription];
    [dataItemEelements_ setObject:value forKey:key];
}

+ (NSString*)getClassDescription:(NSMutableDictionary*)data
{
    NSString* key = [NotificationDataItem getKey:ElementType_ClassDescription];
    return [data objectForKey:key];
}

- (NSString*)getClassMaxStudent
{
    NSString* key = [NotificationDataItem getKey:ElementType_ClassStudent];
    NSString* unique = [dataItemEelements_ objectForKey:key];
    return unique;
}

- (void)setClassMaxStudent:(NSString*)value
{
    NSString* key = [NotificationDataItem getKey:ElementType_ClassStudent];
    [dataItemEelements_ setObject:value forKey:key];
}

+ (NSString*)getClassMaxStudent:(NSMutableDictionary*)data
{
    NSString* key = [NotificationDataItem getKey:ElementType_ClassStudent];
    return [data objectForKey:key];
}


- (void)setNotificationDataItem:(NotificationDataItem*)item
{
    [self setClassName:[item getClassName]];
    [self setClassTime:[item getClassTime]];
    [self setClassDescription:[item getClassDescription]];
    [self setClassTeacher:[item getClassTeacher]];
    [self setClassMaxStudent:[item getClassMaxStudent]];
}

- (NSString*)getUniqueKey
{
    NSString* key = [NotificationDataItem getKey:ElementType_UniqueKey];
    NSString* unique = [dataItemEelements_ objectForKey:key];
    return unique;
}

+ (NSString*)getKey:(ElementType)type
{
    NSString* key;
    switch (type)
    {
        case ElementType_ClassName:
        {
            key = @"ElementType_ClassName";
        }
            break;
        case ElementType_UniqueKey:
        {
            key = @"ElementType_UniqueKey";
        }
            break;
        case ElementType_ClassTime:
        {
            key = @"ElementType_ClassTime";
        }
            break;
        case ElementType_ClassTeacher:
        {
            key = @"ElementType_ClassTeacher";
        }
            break;
        case ElementType_ClassDescription:
        {
            key = @"ElementType_ClassDescription";
        }
            break;
        case ElementType_ClassStudent:
        {
            key = @"ElementType_ClassStudent";
        }
            break;
            
        case ElementType_UserName:
        {
            key = @"ElementType_UserName";
        }
            break;
            
        case ElementType_UserUniqueKey:
        {
            key = @"ElementType_UserUniqueKey";
        }
            break;
        
        case ElementType_ClassRegisteredCount:
        {
            key = @"ElementType_RegisteredStudentCount";
        }
            break;
        
            
        default:
            break;
    }
    return key;
}

+ (NSString*)getCurrentTimeString
{
    NSDateFormatter *dateFormatter = [[NSDateFormatter alloc]init];
    [dateFormatter setDateFormat:@"yyyy-MM-dd HH:mm:ss"];
    
    NSString *dateString = [dateFormatter stringFromDate:[NSDate date]];
    return dateString;
}

+ (NSString*)getUniqueIndex
{
    static int index = 0;
    index++;
    return [NSString stringWithFormat:@"%i", index];;
}

@end

