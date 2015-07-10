//
//  OneClassReigsterStudentData.h
//  Squirrel
//
//  Created by JamesMao on 6/29/15.
//  Copyright (c) 2015 JamesMao. All rights reserved.
//

#import <Foundation/Foundation.h>

@interface OneClassReigsterStudentData : NSObject

+ (OneClassReigsterStudentData*)singleton;

@property (nonatomic, strong)NSString* classUniqueID;

- (int)getUserRegisterDataItemCount;
- (NSMutableDictionary*)getUserRegisterDataItem:(int)index;

- (void)reset;
- (void)save;

- (void)setJsonDictionaryData:(NSDictionary*)jsonDictionary;

@end
