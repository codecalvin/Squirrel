//
//  OneClassReigsterStudentData.m
//  Squirrel
//
//  Created by JamesMao on 6/29/15.
//  Copyright (c) 2015 JamesMao. All rights reserved.
//

#import "OneClassReigsterStudentData.h"

#import "UserRegisterDataItem.h"

#define OneClassReigsterStudentDataFileName @"OneClassReigsterStudentDataFileName.plist"

@interface OneClassReigsterStudentData ()
{
    NSMutableArray* UserRegisterDictionary;
    UserRegisterDataItem* userRegisterDataItem_;
}

@end

@implementation OneClassReigsterStudentData

+ (OneClassReigsterStudentData*)singleton
{
    static OneClassReigsterStudentData* theOneClassReigsterStudentDataInstance;
    if (theOneClassReigsterStudentDataInstance == nil)
    {
        theOneClassReigsterStudentDataInstance =  [[OneClassReigsterStudentData alloc] init];
    }
    return theOneClassReigsterStudentDataInstance;
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
    userRegisterDataItem_ = [[UserRegisterDataItem alloc] init];
    
    NSString* filePath = [self dataFilePath];
    if ([[NSFileManager defaultManager] fileExistsAtPath:filePath])
    {
        UserRegisterDictionary = [[NSMutableArray alloc]initWithContentsOfFile:filePath];
    }
    else
    {
        UserRegisterDictionary = [[NSMutableArray alloc] init];
    }
    
}

- (NSString* )dataFilePath
{
    NSArray* path = NSSearchPathForDirectoriesInDomains(NSDocumentDirectory, NSUserDomainMask, YES);
    NSString* documentsDirectory = [path objectAtIndex:0];
    NSString* dataFilePathFullName = [documentsDirectory stringByAppendingPathComponent:OneClassReigsterStudentDataFileName];
    return dataFilePathFullName;
}

- (int)getUserRegisterDataItemCount
{
    return (int)[UserRegisterDictionary count];
}

- (NSMutableDictionary*)getUserRegisterDataItem:(int)index
{
    if (0 <= index && index < (int)[UserRegisterDictionary count])
    {
        return [UserRegisterDictionary objectAtIndex:index];
    }
    
    return nil;
}

- (BOOL)setUserRegisterDataItem:(UserRegisterDataItem*)dataItem
{
    for (int index = 0; index < [UserRegisterDictionary count]; index++)
    {
        NSMutableDictionary*  dataItemEelements = [self getUserRegisterDataItem:index];
        if (dataItemEelements == nil)
        {
            continue;
        }
        
        NSString* uniqueKeyCandidate = [UserRegisterDataItem getUserUniqueKey:dataItemEelements];
        NSString* uniqueKey = [dataItem getUserUniqueKey];
        if ([uniqueKey compare:uniqueKeyCandidate] == NSOrderedSame)
        {
            [UserRegisterDictionary removeObject:dataItemEelements];
            [UserRegisterDictionary addObject:[NSMutableDictionary dictionaryWithDictionary:[dataItem getDataItemEelements] ]];
            return YES;
        }
    }
    
    NSMutableDictionary* dataItemEelements = [[NSMutableDictionary alloc] initWithDictionary:[dataItem getDataItemEelements]];
    
    [UserRegisterDictionary addObject:dataItemEelements];
    
    [self save];
    return YES;
}

- (void)reset
{
    [UserRegisterDictionary removeAllObjects];
}

- (void)save
{
    NSString* filePath = [self dataFilePath];
    if ([UserRegisterDictionary writeToFile:filePath atomically:YES] == NO)
    {
        [UserRegisterDictionary writeToFile:filePath atomically:NO];
    }
}

- (void)setJsonDictionaryData:(NSDictionary*)jsonDictionary
{
    NSArray* keys = [jsonDictionary allKeys];
    for (NSString* key in keys)
    {
        NSString* value = [jsonDictionary objectForKey:key];
        [userRegisterDataItem_ setUserUniqueKey:key];
        [userRegisterDataItem_ setUserName:value];
        
        [self setUserRegisterDataItem:userRegisterDataItem_];
    }
}


@end

