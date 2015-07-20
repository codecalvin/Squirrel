//
//  LoginViewController.m
//  Squirrel
//
//  Created by JamesMao on 4/11/15.
//  Copyright (c) 2015 JamesMao. All rights reserved.
//

#import "LoginViewController.h"
#import "TabBarViewController/DictionaryTabBarViewController.h"
#import "DataModel/MeViewData.h"
#import "Define.h"
#import "UrlManager/UrlManager.h"

@interface LoginViewController ()
{
    IBOutlet UILabel* errorInformation_;
    
    IBOutlet UIButton* loginButton_;
    IBOutlet UITextField* userNameTextField_;
    IBOutlet UITextField* passwordTextField_;
}
@end

@implementation LoginViewController

- (void)viewDidLoad
{
    [super viewDidLoad];
    
    passwordTextField_.secureTextEntry = YES;
    [self.view setTag:LIGIN_VIEW_TAG];
    errorInformation_.text = @"";
}

-(void)viewWillAppear:(BOOL)animated
{
    [super viewWillAppear:animated];
}

- (id)initWithCoder:(NSCoder *)aDecoder
{
    if ((self = [super initWithCoder:aDecoder])) {
        [[NSBundle mainBundle] loadNibNamed:@"LoginViewController" owner:self options:nil];
    }
    return self;
}

- (void)didReceiveMemoryWarning
{
    [super didReceiveMemoryWarning];
    // Dispose of any resources that can be recreated.
}

- (IBAction)onLoginButton:(id)sender
{
    [self resignKeybard];
    
    if ([self passLoginVerification] == NO)
    {
        return;
    }
    
    [self registerUserInformation];
    
    [self loginViewSwitch];
    //[self loginCheck];
}

- (void) loginCheck {
    NSMutableDictionary *parameters = [[ NSMutableDictionary alloc] init];
    [parameters setValue:userNameTextField_.text forKey:@"ads_name"];
    [parameters setValue:passwordTextField_.text forKey:@"ads_pass"];
    
    NSString* signupURL = [NSString stringWithFormat:@"%@%@", SERVER_IP, URL_SIGNUP_URL];
    [self request:RequestTypePost urlString:signupURL parameters:parameters];
}
- (void)registerUserInformation
{
    if ([userNameTextField_.text compare:[MeViewData singleton].userUniqueName] != NSOrderedSame )
    {
        [[MeViewData singleton] reset];
    }
    
    [MeViewData singleton].userUniqueName = [userNameTextField_.text stringByReplacingOccurrencesOfString:@" " withString:@""];
    [MeViewData singleton].userUniqueName = [[MeViewData singleton].userUniqueName stringByReplacingOccurrencesOfString:@"    " withString:@""];
}

- (BOOL)passLoginVerification
{
    if ([userNameTextField_.text length] > 0 && [passwordTextField_.text length] > 0)
    {
        return YES;
    }
    
    errorInformation_.text = @"User name or password error!";
    

    return NO;
}

- (IBAction)onTouchBackground:(id)sender
{
    [self resignKeybard];
}

- (void)resignKeybard
{
    [userNameTextField_ resignFirstResponder];
    [passwordTextField_ resignFirstResponder];
}

- (void)loginViewSwitch
{
    UIView* window = [self.view superview];
    self.view.hidden = YES;
    [window addSubview:[DictionaryTabBarViewController singleton].view];
    
    passwordTextField_.text = @"";
    
}

- (void)onSuccess:(AFHTTPRequestOperation *)operation responseObject:(id)responseObject {
    [self loginViewSwitch];
}

- (void)onFail:(AFHTTPRequestOperation *)operation error:(NSError*)error {
    errorInformation_.text = @"User name or password error!";
}
@end
