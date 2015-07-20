//
//  URLManager.h
//  Squirrel
//
//  Created by JamesMao on 7/5/15.
//
//

#import <Foundation/Foundation.h>

typedef enum
{
    URLTypeOneClass,
    URLTypeOneClassUserCount,
    URLTypeClassPost,
    URLTypeUserRegister,
    URLTypeUserUnregister,
    URLTypeQueryRegisterStatus,
} URLType;

@interface URLManager : NSObject

+ (NSString*)urlString:(URLType)type;
+ (NSString*)urlString:(URLType)type variableKey:(NSString*)key;

+ (BOOL)isURLString:(URLType)type candidate:(NSString*)candidate;

@end
