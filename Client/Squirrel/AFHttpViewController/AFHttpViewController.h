//
//  AFHttpViewController.h
//  Squirrel
//
//  Created by JamesMao on 7/5/15.
//
//

#import <UIKit/UIKit.h>

@class AFHTTPRequestOperation;

typedef enum
{
    RequestTypeGet,
    RequestTypePost,
    RequestTypeDelete,
} RequestType;

@interface AFHttpViewController : UIViewController

// to be called
- (void)request:(RequestType)type urlString:(NSString*)urlString parameters:(id)parameters;

// to be override
- (void)onSuccess:(AFHTTPRequestOperation *)operation responseObject:(id)responseObject;
- (void)onFail:(AFHTTPRequestOperation *)operation error:(NSError*)error;

@end
