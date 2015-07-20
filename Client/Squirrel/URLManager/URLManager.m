//
//  URLManager.m
//  Squirrel
//
//  Created by JamesMao on 7/5/15.
//
//

#import "URLManager.h"
#import "define.h"

@implementation URLManager

+ (NSString*)urlString:(URLType)type
{
    return [self urlString:type variableKey:nil];
}

+ (NSString*)urlString:(URLType)type variableKey:(NSString*)key
{
    NSString* partURLString = [self partURLString:type variableKey:key];
    NSString* urlString = [NSString stringWithFormat:@"%@%@", SERVER_IP, partURLString];
    return urlString;
}

+ (BOOL)isURLString:(URLType)type candidate:(NSString*)candidate
{
    NSString* urlString = [self urlString:type];
    if ([candidate rangeOfString:urlString].length > 0)
    {
        return YES;
    }
    else
    {
        return NO;
    }
    
}

+ (NSString*)partURLString:(URLType)type
{
    return [self urlString:type variableKey:nil];
}

+ (NSString*)partURLString:(URLType)type variableKey:(NSString*)key
{
    NSString* partString = @"";
    switch (type)
    {
        case URLTypeOneClass:
        {
            partString =  URL_PART_ONE_CLASS;
        }
        break;
        case URLTypeOneClassUserCount:
        {
            partString =  URL_PART_ONE_CLASS_USER_COUNT;
        }
            break;
            
        case URLTypeClassPost:
        {
            partString =  URL_PART_CLASS_POST;
        }
        break;
        case URLTypeUserRegister:
        {
            partString =  URL_PART_REGISTER_ONE_CLASS;
        }
            break;
        case URLTypeQueryRegisterStatus:
        {
            partString =  URL_PART_QUERY_REGISTER_ONE_CLASS;
        }
            break;
        case URLTypeUserUnregister:
        {
            partString =  URL_PART_UNREGISTER_ONE_CLASS;
        }
            break;

        
            
        default:
            break;
    }
    
    if (key != nil)
    {
        partString = [partString stringByAppendingString:key];
    }
    return partString;
}

@end
