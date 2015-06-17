//
//  UserRegisterDataItem.m
//  Squirrel
//
//  Created by JamesMao on 6/16/15.
//  Copyright (c) 2015 JamesMao. All rights reserved.
//

#import "UserRegisterDataItem.h"

@interface UserRegisterDataItem ()
{
    NSMutableDictionary* dataItemEelements_;
}

@end

@implementation UserRegisterDataItem

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
    dataItemEelements_ = data;
}

- (NSMutableDictionary*)getDataItemEelements
{
    return dataItemEelements_;
}


- (NSString*)getUniqueKey
{
    NSString* key = [NotificationDataItem getKey:ElementType_UniqueKey];
    NSString* unique = [dataItemEelements_ objectForKey:key];
    return unique;
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

- (void)setUserUniqueKey:(NSString*)uniqueKey
{
    NSString* key = [NotificationDataItem getKey:ElementType_UserUniqueKey];
    [dataItemEelements_ setObject:uniqueKey forKey:key];
}

- (NSString*)getUserUniqueKey
{
    NSString* key = [NotificationDataItem getKey:ElementType_UserUniqueKey];
    NSString* unique = [dataItemEelements_ objectForKey:key];
    return unique;
}

+ (NSString*)getUserUniqueKey:(NSMutableDictionary*)data
{
    NSString* key = [NotificationDataItem getKey:ElementType_UserUniqueKey];
    return [data objectForKey:key];
}

- (void)setUserName:(NSString*)name
{
    NSString* key = [NotificationDataItem getKey:ElementType_UserName];
    [dataItemEelements_ setObject:name forKey:key];
}

- (NSString*)getUserName
{
    NSString* key = [NotificationDataItem getKey:ElementType_UserName];
    NSString* name = [dataItemEelements_ objectForKey:key];
    return name;
}

+ (NSString*)getUserName:(NSMutableDictionary*)data
{
    NSString* key = [NotificationDataItem getKey:ElementType_UserName];
    return [data objectForKey:key];
}


@end